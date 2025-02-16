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
                `;
                feedContainer.appendChild(feedItem);
            });
        })
        .catch(error => {
            console.error('Error fetching RSS feed:', error);
        });
</script>

[JSON RSS Feed](https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/docs/feed.json)


