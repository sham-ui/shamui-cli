import { default as ShamUI, DI } from 'sham-ui';
import controller from './controllers/main';

DI.bind( 'widget-binder', controller );

if ( module.hot ) {
    const app = DI.resolve( 'widgets:app' );
    if ( app !== undefined ) {
        app.remove();
        DI.resolve( 'sham-ui' ).render.widgets.forEach( widget => {
            try {
                widget.remove();
            } catch ( e ) {
                // eslint-disable-next-line no-empty
            }
        } );
    }
    module.hot.accept();
}

const UI = new ShamUI();

UI.render.FORCE_ALL();

