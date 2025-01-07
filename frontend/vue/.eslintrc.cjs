/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  'extends': [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:vue/vue3-strongly-recommended',
    'plugin:import/warnings',
    'plugin:import/typescript'
  ],
  parserOptions: {
    parser: require.resolve('@typescript-eslint/parser'),
    extraFileExtensions: ['.vue']
  },
  env: {
    browser: true,
    es2021: true,
    node: true,
    'vue/setup-compiler-macros': true
  },
  plugins: [
    '@typescript-eslint',
    'vue'
  ],
  rules: {
    'vue/html-self-closing': ['error', {
      html: {
        void: 'always', // Не использовать автоматическое закрытие для самозакрывающихся тегов (например, <img>)
        normal: 'always', // Не использовать автоматическое закрытие для обычных тегов
        component: 'always' // Всегда использовать автоматическое закрытие для компонентов
      },
      svg: 'always', // Всегда использовать автоматическое закрытие для тегов <svg>
      math: 'always' // Всегда использовать автоматическое закрытие для тегов <math>
    }],

    'vue/html-indent': ['error', 2],

    'vue/max-attributes-per-line': ['error', {
      singleline: 2, // Максимальное количество атрибутов в одной строке
      multiline: {
        max: 1
      }
    }],

    // add space between curly brackets and it's content
    'object-curly-spacing': ['error', 'always'],

    // remove whitespaces at the end of the line
    'no-trailing-spaces': [2, { skipBlankLines: true }],

    // space between block start and word before it: if (a) {
    'space-before-blocks': 'error',

    // space after keyword: if (a)
    'keyword-spacing': ['error', { before: true }],

    // remove trailing comma to objects/arrays/imports/exports at the end
    'comma-dangle': ['error'],

    // no function params reassign,
    'no-param-reassign': ['error'],

    'space-before-function-paren': [
      'error',
      {
        anonymous: 'never',
        named: 'never',
        asyncArrow: 'always'
      }
    ],

    // no space between method call and paren
    '@typescript-eslint/func-call-spacing': ['error', 'never'],

    // Add space after // or /*
    'spaced-comment': ['error', 'always'],

    'arrow-spacing': ['error', { before: true, after: true }],

    'import/newline-after-import': ['error'],
    'import/no-named-as-default': 'off',
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

    'vue/multi-word-component-names': 'off',

    'prefer-promise-reject-errors': 'off',

    // this rule, if on, would require explicit return type on the `render` function
    '@typescript-eslint/explicit-function-return-type': 'off',

    // in plain CommonJS modules, you can't use `import foo = require('foo')` to pass this rule, so it has to be disabled
    '@typescript-eslint/no-var-requires': 'off',

    // use typescript checker for this rule
    'no-undef': 'off',

    // no if(true)
    'no-constant-condition': ['error', { checkLoops: false }],

    // The core 'no-unused-vars' rules (in the eslint:recommended ruleset)
    // does not work with type definitions
    'no-unused-vars': 'off',

    // allow debugger during development only
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',

    // use interface instead of type
    '@typescript-eslint/consistent-type-definitions': ['error', 'interface'],

    // allowed: T[]
    '@typescript-eslint/array-type': ['error', { default: 'array' }],

    '@typescript-eslint/no-non-null-assertion': 'off',

    // prefer const to let
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

    // prefer { foo } to { foo: foo }
    'object-shorthand': ['error', 'always'],

    // forbid bar + 'a'. allow `${bar}a`
    'prefer-template': ['error'],

    // prefer single quotes
    quotes: [
      'error',
      'single',
      {
        avoidEscape: true
      }
    ],

    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: ['variable', 'function'],
        format: ['camelCase'] // allow camelCase variable name
      },
      {
        selector: ['variable'],
        modifiers: ['exported'],
        format: ['camelCase', 'UPPER_CASE'] // allow UPPER_CASE too for consts
      },
      {
        selector: 'typeLike',
        format: ['PascalCase'] // PascalCase for types
      }
    ],

    '@typescript-eslint/no-explicit-any': [
      'off' // must use unknown but not any
    ],
    // must be kebab-case. i.e. emit('some-event')
    'vue/custom-event-name-casing': [
      'error',
      'kebab-case',
      {
        ignores: []
      }
    ],

    'vue/component-tags-order': [
      'error',
      {
        order: ['template', 'script', 'style']
      }
    ],

    '@typescript-eslint/no-unused-vars': 'off',
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
    // убираем ;
    'semi': ['error', 'never'],
    '@typescript-eslint/semi': ['error', 'never'],
    'vue/valid-v-model': ['off'],
    'indent': ['error', 2],

    // убираем лишние пробелы, идущие друг за другом
    'no-mixed-spaces-and-tabs': ['error'],

    // разбиваем построчно условия, добавляем отступы
    'padded-blocks': ['error', 'never'],

    // убираем лишние пустые строки
    'no-multiple-empty-lines': ['error', {
      max: 1,
      maxEOF: 0
    }],

    // убираем лишние пробелы перед двоеточием, добавляем пробелы после двоеточия
    'key-spacing': ['error', { 'beforeColon': false, 'afterColon': true, 'mode': 'strict' }],

    // убираем лишние пробелы перед запятой, добавляем пробелы после запятой
    'comma-spacing': ['error', { 'before': false, 'after': true }],

    // добавляем пробел перед типизацией
    '@typescript-eslint/type-annotation-spacing': 'error',

    // добавляем строчку перед return
    'newline-before-return': 'error',

    // сортировка пропсов компонентов и блоков
    'vue/attributes-order': ['warn', {
      'order': [
        'DEFINITION',
        'LIST_RENDERING',
        'CONDITIONALS',
        'RENDER_MODIFIERS',
        'GLOBAL',
        ['UNIQUE', 'SLOT'],
        'TWO_WAY_BINDING',
        'OTHER_DIRECTIVES',
        'OTHER_ATTR',
        'CONTENT',
        'EVENTS'
      ],
      'alphabetical': true
    }]
  }
}
