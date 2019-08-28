import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'unknown page', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/unknown-page/' );
    await app.start();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
