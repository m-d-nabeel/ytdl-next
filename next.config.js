/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    FLUENTFFMPEG_COV: "",
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "i.ytimg.com",
        port: "",
        pathname: "**",
      },
    ],
  },
  crossOrigin: "anonymous",
};

module.exports = {
  ...nextConfig,
  output: "standalone",
};
