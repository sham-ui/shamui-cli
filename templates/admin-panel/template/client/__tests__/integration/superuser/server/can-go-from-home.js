import setup, { app } from '../../helpers';
import axios from 'axios';

jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'can go from home', async() => {
    expect.assertions( 2 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/server-info', {
            'Host': 'test',
            'Runtime': {
                'NumCPU': 4,
                'Memory': 3312952,
                'MemSys': 71760,
                'HeapAlloc': 3235,
                'HeapSys': 65120,
                'HeapIdle': 60512,
                'HeapInuse': 4608,
                'HeapRealease': 58400,
                'NextGC': 6305,
                'Goroutines': 8,
                'UpTime': 9, 'Time': 'Fri, 30 Oct 2020 21:05:45 +0100'
            },
            'Files': [ {
                'Name': 'dist/bundle.css',
                'Size': 20035,
                'ModTime': 'Thu, 29 Oct 2020 19:30:25 +0100'
            }, {
                'Name': 'dist/bundle.js',
                'Size': 136980,
                'ModTime': 'Thu, 29 Oct 2020 19:30:25 +0100'
            } ]
        } );

    await app.start();
    app.click( '.sideleft .icon-server' );
    await app.waitRendering();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/server' );
} );
