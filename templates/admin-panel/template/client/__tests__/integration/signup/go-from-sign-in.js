import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'can go to sign in page', async() => {
    expect.assertions( 2 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {}, 401 );

    history.pushState( {}, '', 'http://client.example.com/signup/' );

    await app.start();
    app.click( '.signup-label a' );
    await app.waitRendering();

    expect( window.location.href ).toBe( 'http://client.example.com/login' );
    app.checkBody();
} );
