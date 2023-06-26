const path = require('path');

const withPWA = require('next-pwa')({
  disable: process.env.NODE_ENV === 'development',
  dest: 'public',
  register: true,
  skipWaiting: true,
  customWorkerDir: 'src/worker',
});

/**@type {import('next').NextConfig}*/
const nextConfig = {
  /* config options here */
  distDir: 'out',
};

module.exports = withPWA(nextConfig);
