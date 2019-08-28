import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'authenticated member can\'t visit login page', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/login/' );
    await app.start();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
