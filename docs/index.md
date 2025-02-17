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
  <div class="intro">
    <p>Track developer tool releases from <strong>GitHub</strong> and <strong>Vendor Websites</strong> with ease.</p>
    <a href="https://raw.githubusercontent.com/lathanagaraj/versionpulse/refs/heads/main/feed.json">RSS Feed</a>

  </div>
</header>

<section class="scrollable-section">
  <div id="rss-feed-container">
    <div id="rss-feed" class="rss-grid"></div>
  </div>
  <div id="rss-table-container">
    <div id="rss-table" class="table"></div>
  </div>
</section>

<style>
  body {
    font-family: Arial, sans-serif;
    background-color: #f5f6fa;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  /* Fixed Banner */
  .banner {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    background-color: #3b007b;
    color: white;
    padding: 20px;
    text-align: center;
    z-index: 10;
    height: auto;
  }

  .banner h1 {
    margin: 0;
  }

  /* Scrollable Section */
  .scrollable-section {
    margin-top: 120px; /* To make space for the fixed banner */
    padding: 20px;
    overflow-y: auto;
    height: calc(100vh - 120px); /* Allow scrolling for content */
  }

  /* Toggle Button */
  #toggle-view {
    background-color: #3b007b;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 1em;
    margin-bottom: 20px;
  }

  /* Cards Layout */
  #rss-feed-container {
    display: block;
  }



.rss-item {
    flex: 1 1 20%; /* This means each item will take up 20% of the available width */
    max-width: 20%; /* Optional, ensures each item doesn't exceed 20% */
    box-sizing: border-box;
    margin-bottom: 20px;
}
  #rss-feed {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: center;
    align-items: stretch;
  }

  .rss-item {
    flex: 1 1 20%; /* This means each item will take up 20% of the available width */
    max-width: 20%; /* Optional, ensures each item doesn't exceed 20% */
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

  .rss-item h3 {
    font-size: 1.2em;
    margin-bottom: 10px;
    color: #3b007b;
    text-align: center;
  }

  .rss-item p {
    text-align: left;
    margin-bottom: 10px;
  }

  .rss-item h4 {
    text-align: center;
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
    height: auto;
    margin: 0 auto;
    display: none; /* Initially hidden */
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

</style>

<!-- Toggle Button -->
<!-- <button id="toggle-view">Switch to Table View</button> -->

<!-- Cards Section -->
<div id="rss-feed-container">
  <div id="rss-feed"></div>
</div>

<!-- Table Section -->
<div id="rss-table-container">
  <table id="rss-table">
    <thead>
      <tr>
        <th>Title</th>
        <th>Summary</th>
        <th>Content</th>
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

  // Toggle between Card and Table view
//   document.getElementById("toggle-view").addEventListener("click", function() {
//     const cardContainer = document.getElementById("rss-feed-container");
//     const tableContainer = document.getElementById("rss-table-container");

//     if (cardContainer.style.display === "none") {
//       cardContainer.style.display = "block";
//       tableContainer.style.display = "none";
//       this.textContent = "Switch to Table View";
//     } else {
//       cardContainer.style.display = "none";
//       tableContainer.style.display = "block";
//       this.textContent = "Switch to Card View";
//     }
//   });

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
          <h3><a href="${item.url}" target="_blank">${item.title}</a></h3>
          <h4>${item.summary}</h4>
          <p class="content">${item.content_html}</p>
          <p class="published">Published on: ${formatDate(item.date_published)}</p>
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



