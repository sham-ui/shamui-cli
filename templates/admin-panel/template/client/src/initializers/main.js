import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
import Session from '../services/session';
import Store from '../services/store';
import startRouter from './routes';
import App from '../components/App.sfc';

export default function() {
    new Session();
    new Store();

    startRouter();

    new App( {
        ID: 'app',
        container: document.querySelector( 'body' ),
        directives: {
            ...directives,
            hrefto
        }
    } );
}
