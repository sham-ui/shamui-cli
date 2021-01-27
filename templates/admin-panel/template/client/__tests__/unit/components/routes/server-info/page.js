import RoutesServerInfoPage  from '../../../../../src/components/routes/server-info/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        api: {
            request: getMock.mockReturnValueOnce(
                Promise.resolve( { } )
            )
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
        api: {
            request: jest.fn().mockReturnValueOnce(
                Promise.reject( {} )
            )
        }
    } );
    const meta = renderer( RoutesServerInfoPage, {
        directives,
        filters: {}
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
