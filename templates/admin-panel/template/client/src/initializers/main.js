import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import { storage as appState } from '../storages/app';
import Session from '../services/session';
import Store from '../services/store';
import startRouter from './routes';
import App from '../components/App.sfc';

export default function() {

    // Create services
    const session = new Session();
    const store = new Store();

    // Mount root component
    new App( {
        ID: 'app',
        container: document.querySelector( 'body' ),
        directives: {
            ...directives,
            hrefto
        }
    } );

    // Load token
    store.csrftoken().then( () => {
        appState.tokenLoaded = true;
        appState.sync();
        return (

            // Validate session (get session data)
            session.validateSessionPromise && // Operator `&&` - is hack for no-unused-expressions

            // Init router
            startRouter()
        );
    } );
}
