//  @ts-check

import globals from 'globals'
import { tanstackConfig } from '@tanstack/eslint-config'
import pluginReact from 'eslint-plugin-react'
import tailwindcss from 'eslint-plugin-tailwindcss'
import prettier from 'eslint-config-prettier'
import reactHooks from 'eslint-plugin-react-hooks'
import tanstackQuery from '@tanstack/eslint-plugin-query'

export default [
  ...tanstackConfig,
  {
    name: 'ignore-config-files',
    ignores: ['prettier.config.js', 'eslint.config.js', 'vite.config.js'],
  },
  {
    name: 'globals/browser-and-es2024',
    files: ['**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.es2024,
      },
    },
  },
  // pluginReact.configs.flat.recommended,
  {
    name: 'react/recommended',
    ...pluginReact.configs.flat.recommended,
    settings: {
      react: {
        version: 'detect',
      },
    },
  },
  // pluginReact.configs.flat['jsx-runtime'],
  {
    name: 'react/jsx-runtime',
    ...pluginReact.configs.flat['jsx-runtime'],
  },
  // reactHooks.configs.flat.recommended,
  {
    name: 'react-hooks/recommended',
    ...reactHooks.configs.flat.recommended,
  },
  ...tailwindcss.configs['flat/recommended'],
  {
    name: 'tailwindcss/settings',
    rules: {
      'tailwindcss/no-custom-classname': 'off',
    },
    settings: {
      tailwindcss: {
        config: false,
        cssFiles: ['./src/styles.css'],
        callees: ['clsx', 'classnames', 'tw', 'cn'],
      },
    },
  },
  ...tanstackQuery.configs['flat/recommended'],

  // eslint-config-prettier to be last
  // prettier,
  {
    name: 'prettier/disable-formatting-rules',
    ...prettier,
  },
]
