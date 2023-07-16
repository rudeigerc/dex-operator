import { defineConfig } from "vitepress";

const title = "Dex Operator";
const description =
  "Dex Operator is an unofficial Kubernetes operator for Dex, an identity service that uses OpenID Connect to drive authentication for other apps.";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: "/dex-operator/",
  title: title,
  description: description,
  lastUpdated: true,
  cleanUrls: true,
  ignoreDeadLinks: [
    /:\/\/localhost/,
  ],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    logo: "/logo.svg",
    nav: [
      { text: "Home", link: "/" },
      { text: "Introduction", link: "/introduction/getting-started" },
    ],
    footer: {
      message: "Released under the Apache 2.0 License.",
      copyright: "Copyright © 2023 Yuchen Cheng",
    },
    search: {
      provider: "local",
    },
    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "Getting Started", link: "/introduction/getting-started" },
        ],
      },
      {
        text: "Dex Operator",
        items: [{ text: "Design", link: "/operator/design" }],
      },
      {
        text: "Development",
        items: [
          { text: "Contributing", link: "/development/contributing" },
          { text: "Testing", link: "/development/testing" },
          { text: "Release", link: "/development/release" },
        ],
      },
      {
        text: "Reference",
        items: [{ text: "API Reference", link: "/reference/api" }],
      },
    ],
    socialLinks: [
      { icon: "github", link: "https://github.com/rudeigerc/dex-operator" },
    ],
    editLink: {
      pattern: 'https://github.com/rudeigerc/dex-operator/edit/main/docs/:path',
    },
  },
  head: [
    ["link", { rel: "icon", href: "/favicon.svg", type: "image/svg+xml" }],
    [
      "link",
      {
        rel: "alternate icon",
        href: "/favicon.ico",
        type: "image/png",
        sizes: "16x16",
      },
    ],
    ["meta", { name: "author", content: "Yuchen Cheng" }],
    ["meta", { property: "og:type", content: "website" }],
    ["meta", { name: "og:title", content: title }],
    ["meta", { name: "og:description", content: description }],
    ["meta", { name: "twitter:title", content: title }],
    ["meta", { name: "twitter:card", content: "summary_large_image" }],
    ["meta", { name: "twitter:site", content: "@yuchenrcheng" }],
  ],
  markdown: {
    theme: "github-dark-dimmed",
  },
});
