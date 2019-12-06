import { DI, start } from 'sham-ui';
import mainInitializer from './initializers/main';

if ( module.hot ) {
    const store = DI.resolve( 'sham-ui:store' );
    const app = store.findById( 'app' );
    if ( undefined !== app ) {
        app.remove();
        store.forEach( component => {
            try {
                component.remove();
            } catch ( e ) {
                // eslint-disable-next-line no-empty
            }
        } );
    }
    module.hot.accept();
}

mainInitializer();
start();

