module.exports = {
  'env': {
    'browser': true,
    'es2021': true
  },
  'extends': [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react/recommended',
    'plugin:import/recommended',
    'plugin:import/typescript'
  ],
  'overrides': [
    {
      'env': {
        'node': true
      },
      'files': [
        '.eslintrc.{js,cjs}'
      ],
      'parserOptions': {
        'sourceType': 'script'
      }
    }
  ],
  'parser': '@typescript-eslint/parser',
  'parserOptions': {
    'ecmaVersion': 'latest',
    'sourceType': 'module'
  },
  'plugins': [
    '@typescript-eslint',
    'react',
    'import'
  ],
  'rules': {
    'indent': ['error', 2],
    'quotes': ['error', 'single'],
    'semi': ['error', 'never'],
    'no-multiple-empty-lines': ['error', {
      max: 1,
      maxEOF: 0
    }],
    '@typescript-eslint/member-delimiter-style': ['error', {
      multiline: {
        delimiter: 'semi',
        requireLast: true
      },
      singleline: {
        delimiter: 'semi',
        requireLast: false
      }
    }],
    'padded-blocks': ['error', 'never'],
    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: ['variable', 'function'],
        format: ['camelCase', 'PascalCase', 'snake_case']
      },
      {
        selector: ['variable'],
        modifiers: ['exported'],
        format: ['camelCase', 'UPPER_CASE', 'PascalCase']
      },
      {
        selector: 'typeLike',
        format: ['PascalCase']
      }
    ],
    'object-shorthand': ['error', 'always'],
    'prefer-const': ['error'],
    'prefer-destructuring': [
      'error',
      {
        VariableDeclarator: {
          array: false, // allow: const bar = array[0]
          object: true // disallow: const bar = array.bar
        },
        AssignmentExpression: {
          array: false, // allow: bar = array[0]
          object: false // allow: bar = array.bar
        }
      },
      {
        enforceForRenamedProperties: false // allow bar = foo.baz
      }
    ],
    'import/no-named-as-default': 'off', // почему-то не работает
    'import/newline-after-import': ['error'],
    'import/order': [
      'error',
      {
        groups: ['builtin', 'external', 'internal', 'parent', 'sibling', 'index', 'unknown'],
        pathGroups: [],
        'newlines-between': 'never', // no gaps between imports
        pathGroupsExcludedImportTypes: [],
        alphabetize: {
          order: 'asc',
          caseInsensitive: true
        }
      }
    ],
    'object-curly-spacing': ['error', 'always'],
    'no-trailing-spaces': [2, { skipBlankLines: true }],
    'space-before-blocks': 'error',
    'keyword-spacing': ['error', { before: true }],
    'comma-dangle': ['error'],
    'spaced-comment': ['error', 'always'],
    'arrow-spacing': ['error', { before: true, after: true }],
    'react/function-component-definition': [
      2,
      {
        namedComponents: 'arrow-function',
        unnamedComponents: 'arrow-function'
      }
    ],
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/no-explicit-any': ['off']
  }
}
