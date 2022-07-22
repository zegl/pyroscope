const path = require('path');

module.exports = {
  extends: [path.join(__dirname, '../../../.eslintrc.js')],
  parserOptions: {
    tsconfigRootDir: __dirname,
  },
  overrides: [
    {
      files: ['*.tsx', '*.ts', '*.js'],
      rules: {
        '@typescript-eslint/no-unsafe-assignment': 'off',
        '@typescript-eslint/no-unsafe-call': 'off',
        'react/jsx-props-no-spreading': 'off',
        'no-underscore-dangle': 'off',
        'react/static-property-placement': 'off',
        '@typescript-eslint/ban-types': 'off',
        'no-plusplus': 'off',
        '@typescript-eslint/no-explicit-any': 'off',
        camelcase: 'off',
        'react/state-in-constructor': 'off',
        'no-continue': 'off',
        'react/sort-comp': 'off',
        'jsx-a11y/anchor-is-valid': 'off',
      },
    },
  ],
};
