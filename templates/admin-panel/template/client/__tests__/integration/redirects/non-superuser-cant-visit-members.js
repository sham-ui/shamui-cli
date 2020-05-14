import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'members', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/members/' );
    await app.start();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
