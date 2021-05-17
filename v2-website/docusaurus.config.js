/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  "title": "Magma Documentation",
  "tagline": "Bring more people online by enabling operators with open, flexible, and extensible network solutions",
  "url": "https://magmacore.org",
  "baseUrl": "/",
  "organizationName": "magma",
  "projectName": "magma",
  "scripts": [
    "https://buttons.github.io/buttons.js"
  ],
  "favicon": "img/icon.png",
  "customFields": {
    "disableTitleTagline": true
  },
  "onBrokenLinks": "log",
  "onBrokenMarkdownLinks": "log",
  "presets": [
    [
      "@docusaurus/preset-classic",
      {
        "docs": {
          "showLastUpdateAuthor": true,
          "showLastUpdateTime": true,
          "editUrl": "https://github.com/magma/magma/tree/master/docs",
          "path": "./docs", // FS content path (./docs)
          "routeBasePath": "docs", // Url subpath (http://localhost:3000/docs)
          "sidebarPath": require.resolve("./sidebars.json"),
          // "onlyIncludeVersions": ['current','1.4.0','1.3.0'],
        },
        "blog": {},
        "theme": {
          "customCss": "../src/css/customTheme.css"
        }
      }
    ]
  ],
  "plugins": [],
  "themeConfig": {
    "navbar": {
      "title": "Magma Documentation",
      "logo": {
        "src": "img/magma-logo-purple.svg"
      },
      "items": [
        {
          "href": "https://magmacore.org",
          "label": "Home",
          "position": "left"
        },
        {
          "to": "/docs/basics/introduction",
          "label": "Docs",
          "position": "left"
        },
        {
          type: 'docsVersionDropdown',
          position: "left",
        },
        {
          "href": "https://github.com/magma",
          "label": "Code",
          "position": "left"
        },
        {
          "href": "https://magmacore.org/community",
          "label": "Community",
          "position": "left"
        }
      ]
    },
    "image": "img/docusaurus.png",
    "footer": {
      "links": [],
      "copyright": "Copyright © 2021 The Magma Authors",
      "logo": {
        "src": "img/magma_icon.png"
      }
    },
    "algolia": {
      "apiKey": "f95caeb7bc059b294eec88e340e5445b",
      "indexName": "magma"
    }
  }
});
