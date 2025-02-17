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

    /* Cards Layout */
    #rss-feed {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;
        justify-content: center;
        align-items: stretch;
        margin-bottom: 40px;
    }

    .rss-item {
        width: 320px;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        border: 1px solid #ddd;
        padding: 15px;
        border-radius: 10px;
        box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.1);
        background-color: white;
    }

    .rss-item-container {
        display: flex;
        flex-direction: column;
        height: 100%;
    }

    .rss-item h3 {
        font-size: 1.2em;
        margin-bottom: 10px;
        color: #3b007b;
        text-align: left;
    }

    .rss-item p {
        text-align: left;
        margin-bottom: 10px;
    }

    .rss-item .content {
        flex-grow: 1;
        text-align: left;
        color: #555;
    }

    .rss-item .published {
        font-size: 0.9em;
        color: #888;
        margin-top: auto;
        padding-top: 10px;
        border-top: 1px solid #eee;
    }

    .rss-item a {
        text-decoration: none;
        font-weight: bold;
        color: #3b007b;
    }

    /* Table Layout */
    #rss-table-container {
        width: 90%;
        margin: 0 auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        background: white;
        border-radius: 10px;
        overflow: hidden;
        box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.1);
    }

    th, td {
        border: 1px solid #ddd;
        padding: 10px;
        text-align: left;
    }

    th {
        background-color: #3b007b;
        color: white;
    }

    tr:nth-child(even) {
        background-color: #f2f2f2;
    }

    .table-title {
        text-align: center;
        font-size: 1.5em;
        margin-bottom: 10px;
        color: #3b007b;
    }
</style>

<!-- Cards Section -->
<h2 style="text-align: center; color: #3b007b;">Feed in Card Format</h2>
<div id="rss-feed"></div>

<!-- Table Section -->
<h2 style="text-align: center; color: #3b007b;">Feed in Table Format</h2>
<div id="rss-table-container">
    <table id="rss-table">
        <thead>
            <tr>
                <th>Title</th>
                <th>Summary</th>
                <th>Content</th>
                <th>Published Date</th>
            </tr>
        </thead>
        <tbody></tbody>
    </table>
</div>

<script>
    // Function to format the published date into a readable format
    function formatDate(isoString) {
        const date = new Date(isoString);
        return date.toLocaleString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: '2-digit',
            hour12: true
        });
    }

    // Replace with your RSS feed URL
    const rssUrl = 'https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/feed.json';

    // Fetch RSS feed data and display it
    fetch(rssUrl)
        .then(response => response.json())
        .then(data => {
            const feedContainer = document.getElementById('rss-feed');
            const tableBody = document.querySelector("#rss-table tbody");

            let maxHeight = 0; // Track max height

            data.items.forEach(item => {
                // Create card layout
                const feedItem = document.createElement('div');
                feedItem.classList.add('rss-item');

                feedItem.innerHTML = `
                    <div class="rss-item-container">
                        <h3><a href="${item.url}" target="_blank">${item.title}</a></h3>
                        <p>${item.summary}</p>
                        <p class="content">${item.content_html}</p>
                        <p class="published">Published on: ${formatDate(item.date_published)}</p>
                    </div>
                `;
                feedContainer.appendChild(feedItem);

                // Update maxHeight based on the tallest card
                if (feedItem.clientHeight > maxHeight) {
                    maxHeight = feedItem.clientHeight;
                }

                // Create table row
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td><a href="${item.url}" target="_blank">${item.title}</a></td>
                    <td>${item.summary}</td>
                    <td>${item.content_html}</td>
                    <td>${formatDate(item.date_published)}</td>
                `;
                tableBody.appendChild(row);
            });

            // Apply maxHeight to all cards
            document.querySelectorAll('.rss-item').forEach(item => {
                item.style.height = maxHeight + 'px';
            });
        })
        .catch(error => {
            console.error('Error fetching RSS feed:', error);
        });
</script>

[JSON RSS Feed](https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/feed.json)


