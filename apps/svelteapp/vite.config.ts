import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import { nxViteTsPaths } from '@nx/vite/plugins/nx-tsconfig-paths.plugin';

export default defineConfig({
  root: __dirname,
  cacheDir: '../../node_modules/.vite/svelteapp',
  server: {
    fs: {
      // Allow serving files from one level up to the project root
      allow: ['../..']
    }
  },
  plugins: [sveltekit(), nxViteTsPaths()],
  test: {
    cache: {
      dir: '../../node_modules/.vitest'
    },
    include: ['src/**/*.{test,spec}.{js,ts}'],
    coverage: {
      reportsDirectory: '../../coverage/libs/auth',
      provider: 'v8'
    }
  }
});
