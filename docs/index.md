---
layout: home
author_profile: true
title: "VersionPulse"
description: "Track the latest developer tool releases with RSS feeds."
header:
  caption: "Stay ahead with real-time updates"
excerpt: "VersionPulse aggregates GitHub and vendor releases into a single RSS feed."
---
<head>
    <link rel="stylesheet" href="{{ site.baseurl }}/assets/style.css">
</head>
  

<header class="banner">
  <h1>VersionPulse</h1>
</header>

<div class="intro">
  <p>VersionPulse is your go-to solution for tracking developer tool releases from <strong>GitHub</strong> and <strong>Vendor Websites</strong>. Stay informed with <strong>automated RSS feeds</strong> and never miss an update.</p>
</div>

<section class="scrollable-section">
  <h2>Latest Releases</h2>
  <div id="rss-feed" class="rss-grid"></div>
</section>

<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f5f6fa;
    }

    #rss-feed {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;
        justify-content: center;
    }

    .rss-item {
        width: 320px;
        height: 420px;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        border: 1px solid #ddd;
        padding: 15px;
        border-radius: 10px;
        box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.1);
        background-color: white;
        overflow: hidden;
    }

    .rss-item h3 {
        font-size: 1.2em;
        margin-bottom: 10px;
        color: #3b007b;
    }

    .rss-item p {
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 3; /* Limits summary to 3 lines */
        -webkit-box-orient: vertical;
        margin-bottom: 10px;
    }

    .rss-item .content {
        flex-grow: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 3; /* Limits content to 3 lines */
        -webkit-box-orient: vertical;
        color: #555;
    }

    .rss-item .published {
        font-size: 0.9em;
        color: #888;
        margin-top: auto; /* Ensures it stays at the bottom */
        padding-top: 10px;
        border-top: 1px solid #eee;
    }

    .rss-item a {
        text-decoration: none;
        font-weight: bold;
        color: #3b007b;
    }
</style>

<div id="rss-feed"></div>

<script>
    // Replace with your RSS feed URL
    const rssUrl = 'https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/docs/feed.json';

    // Fetch RSS feed data and display it
    fetch(rssUrl)
        .then(response => response.json())
        .then(data => {
            const feedContainer = document.getElementById('rss-feed');
            data.items.forEach(item => {
                const feedItem = document.createElement('div');
                feedItem.classList.add('rss-item');
                feedItem.innerHTML = `
                    <h3><a href="${item.url}" target="_blank">${item.title}</a></h3>
                    <p>${item.summary}</p>
                    <p class="content">${item.content_html}</p>
                    <p class="published">Published on: ${item.date_published}</p>
                `;
                feedContainer.appendChild(feedItem);
            });
        })
        .catch(error => {
            console.error('Error fetching RSS feed:', error);
        });
</script>



[JSON RSS Feed](https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/docs/feed.json)


