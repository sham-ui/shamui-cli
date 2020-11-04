import RoutesServerInfoPage  from '../../../../../src/components/routes/server-info/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        axios: {
            get: getMock.mockReturnValueOnce(
                Promise.resolve( { } )
            )
        },
        constructor: {
            extractData: x => x
        }
    } );

    const meta = renderer( RoutesServerInfoPage, {
        directives,
        filters: {}
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    DI.resolve( 'session:storage' ).sessionValidated = true;
    DI.bind( 'store', {
        axios: {
            get: jest.fn().mockReturnValueOnce(
                Promise.reject( {} )
            )
        },
        constructor: {
            extractData: x => x
        }
    } );
    const meta = renderer( RoutesServerInfoPage, {
        directives,
        filters: {}
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
