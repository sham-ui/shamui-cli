import { start, DI } from 'sham-ui';
import pretty from 'pretty';
import initializer from '../../src/initializers/main';

export const app = {
    async start( waitRendering = true ) {
        initializer();
        start();
        if ( waitRendering ) {
            await this.waitRendering();
        }
    },
    async waitRendering() {
        await new Promise( resolve => setImmediate( resolve ) );
    },
    click( selector ) {
        document
            .querySelector( selector )
            .click();
    },
    checkBody() {
        expect(
            pretty( document.querySelector( 'body' ).innerHTML, {
                inline: [ 'code', 'pre', 'em', 'strong', 'span' ]
            } ),
        ).toMatchSnapshot();
    },
    form: {
        fill( field, value ) {
            document.querySelector( `[name="${field}"]` ).value = value;
        },
        async submit() {
            app.click( '[type="submit"]' );
            await app.waitRendering();
        }
    }
};

function setupRAF() {
    window.requestAnimationFrame = ( cb ) => {
        setImmediate( cb );
    };
}

function clearBody() {
    document.querySelector( 'body' ).innerHTML = '';
}

function resetShamUI() {
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
}

function resetStorage() {
    DI.resolve( 'session:storage' ).reset();
}

function setupRouter() {
    delete window.__NAVIGO_WINDOW_LOCATION_MOCK__;
    history.pushState( {}, '', '' );
}

export default function() {
    setupRAF();
    resetShamUI();
    clearBody();
    resetStorage();
    setupRouter();
    Object.defineProperty( window, 'CSS', { value: () => ( {} ) } );
}
