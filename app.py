import scrapy
import re
from scrapy.crawler import CrawlerProcess
from scrapy.linkextractors import LinkExtractor
from scrapy.spiders import CrawlSpider, Rule
from scrapy_playwright.page import PageMethod

class GenericSpider(CrawlSpider):
    name = "generic"
    start_urls = ["https://blog.jetbrains.com/idea/category/releases/"]  # Replace with your target URL

    # Define a deny pattern to filter download links
    deny_pattern = re.compile(r".*\.(exe|zip|msi|tar\.gz|dmg|deb|rpm|apk|pkg|iso)$", re.IGNORECASE)

    # Use Scrapy's built-in deny_extensions to prevent requests to unwanted file types
    rules = (
        Rule(
            LinkExtractor(
#                deny_extensions=["exe", "zip", "msi", "tar.gz", "dmg", "deb", "rpm", "apk", "pkg", "iso"]
            ),
            callback="parse_item",
            follow=True,
        ),
    )

    def parse_item(self, response):

         # Extract all visible text
        all_text = " ".join(
            text.strip()
            for text in response.css("body *::text").getall()
            if text.strip() and not self.is_excluded_element(response)
        )

        # Extract JavaScript-related warnings like "Learn how to turn on JavaScript"
        js_warning = response.xpath("//*[contains(text(), 'Learn how to turn on JavaScript')]/text()").get()

       # Extract headers (h1 to h6)
        headers = response.css("h1::text, h2::text, h3::text, h4::text, h5::text, h6::text").getall()

        # Extract all paragraph text
        paragraphs = response.css("p::text").getall()

        # Extract all links
        links = response.css("a::attr(href)").getall()

        # Filtering links to avoid download links
        # links_to_follow = [link for link in links if not re.match(self.deny_pattern, link)]

        yield {
            "url": response.url,
            "title": response.css("title::text").get(),
            "headers": headers,  # Structured headers content
            "paragraphs": paragraphs,  # Structured paragraphs
            "links": links,  # Valid links (no downloads)
            "js_warning": js_warning,  # JavaScript warning content
            "all_text": all_text,  # Cleaned-up text
        }
    def start_requests(self):
        for url in self.start_urls:
            yield scrapy.Request(
                url,
                callback=self.parse_item,
                meta={
                "playwright": True,
                "playwright_page_methods": [
                    PageMethod("wait_for_selector", "body"),  # Wait for the page to load
                ],
                "playwright_include_page": True,},            
            )
    def is_excluded_element(self, response):
        """Checks if the element is a script or style tag or contains inline JS"""
        return response.css("script") or response.css("style")
    

def main():
    process = CrawlerProcess(settings={
        "DOWNLOAD_HANDLERS": {
            "http": "scrapy_playwright.handler.ScrapyPlaywrightDownloadHandler",
            "https": "scrapy_playwright.handler.ScrapyPlaywrightDownloadHandler",
        },
        "TWISTED_REACTOR": "twisted.internet.asyncioreactor.AsyncioSelectorReactor",
        "FEEDS": {
            "output.json": {"format": "json"},
        },
        "DEPTH_LIMIT": 5,  # Limit the depth of crawling (optional)
        "PLAYWRIGHT_BROWSER_TYPE": "chromium",  # Use Chromium for rendering
        "PLAYWRIGHT_PROCESS_REQUEST_HEADERS": None,  # Prevents Playwright header conflicts
        #"DOWNLOAD_MAXSIZE": 1048576,  # 1MB max size for safety
        #"DOWNLOAD_WARNSIZE": 524288,  # 512KB warn size
    })
    process.crawl(GenericSpider)
    process.start()


if __name__ == "__main__":
    main()
