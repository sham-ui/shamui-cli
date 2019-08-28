import ShamUI, { DI } from 'sham-ui';
import pretty from 'pretty';
import initializer from '../../src/initializers/main';

export const app = {
    async start( waitRendering = true ) {
        DI.bind( 'component-binder', initializer );
        new ShamUI( true );
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
    const UI = DI.resolve( 'sham-ui' );
    if ( undefined !== UI ) {
        UI.render.unregister( 'app' );
        DI.resolve( 'sham-ui:store' ).forEach( component => {
            try {
                UI.render.unregister( component.ID );
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
