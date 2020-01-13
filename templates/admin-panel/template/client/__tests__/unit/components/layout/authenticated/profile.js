import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/href-to';
// eslint-disable-next-line max-len
import LayoutAuthenticatedProfile  from '../../../../../src/components/layout/authenticated/profile.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

afterEach( () => {
    DI.resolve( 'session:storage' ).reset();
    DI.bind( 'router', null );
} );

it( 'renders correctly', () => {
    DI.resolve( 'session:storage' ).name = 'Test member';
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );
    const meta = renderer( LayoutAuthenticatedProfile, {
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
