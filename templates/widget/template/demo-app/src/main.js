import ShamUI, { DI } from 'sham-ui';
import controller from './controllers/main';

DI.bind( 'widget-binder', controller );

// Used only for dev, don't include in production build
if ( !PRODUCTION && module.hot ) {
    const UI = DI.resolve( 'sham-ui' );
    if ( undefined !== UI ) {
        UI.render.unregister( 'app' );
        DI.resolve( 'sham-ui:store' ).forEach( widget => {
            try {
                UI.render.unregister( widget.ID );
            } catch ( e ) {
                // eslint-disable-next-line no-empty
            }
        } );
    }
    module.hot.accept();
}

new ShamUI( true );

