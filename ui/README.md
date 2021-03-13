## Info

- Dodałem przykładowy feedback walidacji (elementy z klasą .invalid-feedback lub .valid-feedback)

- FAQ: Jezeli lista FAQ bedzie krotsza niz odpowiedz po prawej, tomoze wystapic "migniecie", to przez JS zapewniajacy odpowiednia wysokosc elementów, trzeba rozciagnac recznie JSem w przypadku wyzszej odpowiedzi (aby pozostawic to jak najprostsze w strukturze/css)

- Umieszczałem grafiki inline (same SVG), ale jeżeli trzeba, to przeniosę do osobnego folderu.

- Email: skorzystałem z https://editor.bootstrapemail.com/ do wygenerowania szablonu pliki wejściowe i wyjściowe są w src/assets). Wrzuciłem PNG inline z src w base64, żeby skrzynki pocztowe to obsługiwały trzeba zamienić na rzeczywiste linki. Można tez wrzucać SVG inline, ale np. gmail desktop tego nie obsługuje.

    Jest:
    ```
    <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAA..." />
    ```
    Powinno być:
    ```
    <img src="https://link.com/image.png" />
    ```

    Pliki z których wygenerowałem data w base64 są w src/assets

### JS Validation

REQUIRED STRUCTURE FOR VALIDATION:
form .row .col .form-group input.form-control

# Symbiosis UI (based on Webpack Frontend Starterkit and Bootstrap npm starter template)

https://github.com/twbs/bootstrap-npm-starter
https://github.com/wbkd/webpack-starter/tree/main/webpack

### Installation

```
npm install
```

### Start Dev Server

```
npm start
```

### Build Prod Version

```
npm run build
```

### Features:

- ES6 Support via [babel](https://babeljs.io/) (v7)
- JavaScript Linting via [eslint-loader](https://github.com/MoOx/eslint-loader)
- SASS Support via [sass-loader](https://github.com/jtangelder/sass-loader)
- Autoprefixing of browserspecific CSS rules via [postcss](https://postcss.org/) and [autoprefixer](https://github.com/postcss/autoprefixer)
- Style Linting via [stylelint](https://stylelint.io/)

When you run `npm run build` we use the [mini-css-extract-plugin](https://github.com/webpack-contrib/mini-css-extract-plugin) to move the css to a separate file. The css file gets included in the head of the `index.html`.

## Other Scripts

The following npm scripts are available to you in this starter repo. With the exception of `npm start`, the remaining scripts can be run from your command line with `npm run scriptName`.

| Script | Description |
| --- | --- |
| `serve` | Starts a local server (<http://localhost:3000>) with /public hosted |

### Optimizing JS

Similar to optimizing CSS, we publish individual scripts for each of our plugins. This allows you to import only what you need, versus the entire bundle and dependencies. For example, if you don't plan on using dropdowns, tooltips, or popovers, you can safely omit the Popper.js depdendency. Bootstrap 4 requires jQuery though, so you won't be able to safely remove that until v5 launches.

See the `src/js/app.js` file for an example of how to import all of Bootstrap's JS or just the individual pieces. By default we've only imported our modal JavaScript since we have no need for anything else.

You can add more options here, or import the entire `bootstrap-bundle.min.js` file, to get all JavaScript plugins and Popper.js.

## Actions CI

We've included some simple GitHub Actions in this template repo. When you generate your new project from here, you'll have the same tests that run whenever a pull request is created. We've included Actions for the following:

- Stylelint for your CSS

When your repository is generated, you won't see anything in the Actions tab until you create a new pull request. You can customize these Actions, add new ones, or remove them outright if you wish.

[Learn more about GitHub Actions](https://github.com/features/actions), [read the Actions docs](https://help.github.com/en/actions), or [browse the Actions Marketplace](https://github.com/marketplace/actions).

## Copyright

&copy; @mdo 2020 and licensed MIT.
