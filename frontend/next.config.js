/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
  redirects: async () => {
    return [
      {
        source: "/",
        destination: "/program/search",
        permanent: false,
      },
    ]
  },
}

module.exports = nextConfig
