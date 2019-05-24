# FBC I18n

This package utilizes [i18next](https://www.i18next.com) and [react-i18next](https://react.i18next.com).

Localization is broken up into 3 key parts:

* Clientside substitution
* Phrase extraction
* Language detection and translation retrieval


### Language Detection
Language detection is the process of determining the user's language and serving them the proper translation files. Whenever a user loads the page, the backend will try to detect their current language. First it checks the querystring for `?locale=<locale>` - if this is set, it overrides all others. Next, it checks the url path for `/:locale`. Finally, it checks if there is a cookie which caches the user's language. If the language is detected through either of the first 2 methods, it is cached in this cookie. That way `?locale=en` doesn't have to remain in their querystring forever. To change the user's language, simply redirect them with the querystring `?locale=<locale>`.

Once the user's language is detected, the translations file is served to the client where it can be used to perform substitution in the client.

### Clientside Substitution
The basic process of clientside substitution is:
* Create a translations file for each locale (en, es, fr, ru)
* For each piece of text in the application, create a key value pair in english
  * en.json `{ "hello":"hello", "thankyou":"thank you" }`
  * es.json `{ "hello":"hola", "thankyou":"gracias" }`
* Find the json file which corresponds to the user's selected locale. If the file is not available, load a fallback such as english
* Instead of rendering the literal text "hello", the app requests the _translation_ for the key hello. In spanish the translation engine will return hola, and in english it will still return hello.


### Phrase Extraction
_This process only occurs during development_

Normally, when building localized applications, the developer must maintain a localized translations file such as `locales/en.json`. Whenever new text is added to the application, it must be added to this file manually. Phrase extraction automates this process. Instead of maintaining the file manually, the developer simply provides a key and the english translation directly inside the source code. This is then extracted at runtime and automatically added to the default translations file. For example:
```js
// Trans is the react component used to handle translations
<Trans i18nKey="extract_me">
  here is some text that needs to be <strong>extracted</strong>!
</Trans>
```

If the translation for the key extract_me is not loaded when this component renders, the key and the english text will be posted to the backend.

More info on translation extraction:
* [Missing Keys](https://www.i18next.com/overview/configuration-options#missing-keys)
* [Runtime Extraction](https://react.i18next.com/guides/extracting-translations#3-runtime-extraction)



## Setup and Examples:

Setup will follow the same 3 key parts:
* Language Detection
* Clientside Substitution
* Phrase Extraction

### Language Detection
* add @fbcnms-i18n to package.json
* create a file like i18n.js *serverside* - add the following code:

```js
import {i18nBuilder, makeLocaleDirectory} from '@fbcnms/i18n';
export const i18nextInstance = i18nBuilder({
  preload: ['en_US'], // which languages to load at startup
  fallbackLng: 'en_US', // which language to fallback to if the requested is not available
});
```

This will serve as the initialized singleton instance of i18next.

i18nextBuilder takes one parameter which configures i18next's [options object](https://www.i18next.com/overview/configuration-options). If an object is passed, it is shallow merged into the default config provided by i18nBuilder, if a configuration function is provided, the default config is passed to the function, this return value becomes the options object.

_Note that i18nBuilder is just a function which provides some sensible i18next defaults. It's possible to completely build a custom instance of i18next. This is useful if, for example, you don't want to store translations on the filesystem, or you don't like the language detection mechanisms._

Most apps have a catch-all route which serves the default HTML page. I recommend detecting the user's language in this route and storing it in the served HTML page, then reading it from JS. Here's an example:

```js
app.get('*', (req, res) => {
  const detectedLanguage = i18nextInstance.services.languageDetector.detect(
    req,
    res,
  );
  return res.render('index', {
    data: { // this object will appear as window.CONFIG in the clientside examples
      LANG:detectedLanguage,
      DEVELOPMENT: process.env.NODE_ENV !== 'production'
    }
  });
})
```

This example handles detecting the user's language via the querystring, caching it in a cookie, and passing it to the rendered html page.


### Clientside Substitution

* `yarn add i18next react-i18next i18next-xhr-backend`
* create a file like i18n.js *clientside* - add the following code:

```js
import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import XHR from 'i18next-xhr-backend';

// certain actions should only occur in dev mode, detect this however makes sense for you.
const DEVELOPMENT = window.CONFIG.env.NODE_ENV !== 'production';
const locale = window.CONFIG.LANG;
i18n
  .use(XHR)
  .use(initReactI18next)
  .init({
    lng: locale,
    fallbackLng: 'en_US',
    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
    backend: {
      // api routes to load translations from, and to save missing translations
      loadPath: '/translations/{{lng}}/{{ns}}.json',
      addPath: '/translations/add/{{lng}}/{{ns}}',
    },
    keySeparator: false,
    debug: DEVELOPMENT,
    saveMissing: DEVELOPMENT,
  });

export default i18n;
```

Make sure to import i18n.js in one of the topmost files of the app, this will trigger translation loading.

In one of your UI files, add the following code:

```js
import Text from '@fbcnms/i18n';
import {withTranslation, useTranslation, Trans} from 'react-i18next';

function App() {
  return (
    <Suspense fallback={<span>Loading...</span>}>
      <MyComponent />
      <MyOtherComponent />
    </Suspense>
  );
}

function MyComponent() {
  // useTranslation triggers a suspense if translations aren't loaded yet
  const {t} = useTranslation()
  // functional style - useful for building strings outside of the render method
  const world_translated = t('world', 'world');
  return (
    <div>
      {/*component style - useful for building strings inside the render method. Also supports arbitrary html in the phrases. Arbitrary html is replaced with tokens like <0></0> and then re-substituted back in clientside. This is useful for styling pieces of text. */}
      <Trans i18nKey="hello">Hello!</Trans>
      <span>{world_translated}</span>
      { /** Text is a component which combines Trans and material-ui's
      Typography component, use it to replace the usage of Typography in your app.
      Typography still has its place for 100% user generated strings*/}
      <Text i18nKey="material_typography" variant="body">material ui typography!</Text>
    </div>
  );
}

// withTranslation triggers a suspense if translations aren't loaded yet
const MyOtherComponent = withTranslation()(function MyOtherComponent({t}){
  return (
    <div>{t('hoc_test', 'HOC style!')}</div>
  )
})
```

_Suspense_

TLDR: If using the XHR backend, make at least one call to useTranslation or withTranslation high in the component tree to prevent flashes of fallback language -> selected language.

In the comments above, I mention [React Suspense](https://reactjs.org/docs/react-api.html#reactsuspense). useTranslation and withTranslation will both force the app to show a loading spinner until translations are loaded. No matter how many calls are made to useTranslation / withTranslation, suspense is only triggered once so performance is not impacted. Another way around this is to remove the xhr backend and bootstrap translations in the HTML. This is easily supported, simply provide the parsed JSON i18n.init

```js
i18n.init({
  resources: {
      en: {
        // default namespace
        translation: {
          hello:'hello'
        }
      }
  }
})
```


__Loading Translations__
* Note that the following steps only apply when using i18next-xhr-backend to load translations.

Now that we've setup the client, we need to create the api routes to load translations from.

When configuring the clientside i18next backend, we provided the following options:

```js
backend: {
  // api routes to load translations from, and to save missing translations
  loadPath: '/translations/{{lng}}/{{ns}}.json',
  addPath: '/translations/add/{{lng}}/{{ns}}',
},
```

Now we must create these routes in the backend

Create the following express route:

```js
import {getLocaleHandler} from '@fbcnms/i18n';
 // we created the file i18n during the language detection step
import {i18nextInstance} from 'i18n';

app.use('/translations/:locale/:ns.json', getLocaleHandler(i18nextInstance));
```

This is all that's required!


### Phrase Extraction

To enable runtime phrase extraction, just add the following route

```js
import {addMissingKeysHandler} from '@fbcnms/i18n';
router.post('/translations/add/:locale/:ns', addMissingKeysHandler(i18nextInstance));
```

**Note that this route will only accept requests in development mode (NODE_ENV !== 'production')**



## Further Reading
* https://www.i18next.com - i18next supports many advanced translation features. I've only demonstrated basic string replacement translations but i18next supports gender, plurals, dynamic interpolation, etc.
* https://react.i18next.com
