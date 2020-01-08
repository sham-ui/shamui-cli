import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import Session from '../services/session';
import Store from '../services/store';
import startRouter from './routes';
import App from '../components/App.sht';

export default function() {
    new Session();
    const store = new Store();

    startRouter();

    const app = new App( {
        ID: 'app',
        container: document.querySelector( 'body' ),
        directives: {
            ...directives,
            hrefto
        },
        tokenLoaded: false,
        routerResolved: false
    } );

    store.tokenPromise.then(
        () => app.update( {
            tokenLoaded: true
        } )
    );
}
