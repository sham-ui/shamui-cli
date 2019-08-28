import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
    window.matchMedia = jest.fn().mockImplementation(
        () => ( {
            addListener: jest.fn(),
            matches: true
        } )
    );
} );

afterEach( () => {
    delete window.matchMedia;
} );

it( 'profile dropdown', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();
    await app.start();
    app.click( '.link-profile' );
    app.checkBody();
    app.click( '.dropdown-menu' );
    app.checkBody();
} );
