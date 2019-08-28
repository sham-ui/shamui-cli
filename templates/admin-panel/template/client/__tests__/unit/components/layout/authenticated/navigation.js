import { DI } from 'sham-ui';
import hrefto from 'sham-ui-router/href-to';
// eslint-disable-next-line max-len
import LayoutAuthenticatedNavigation  from '../../../../../src/components/layout/authenticated/navigation.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

afterEach( () => {
    DI.bind( 'router', null );
} );

it( 'renders correctly', () => {
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        lastRouteResolved: jest.fn().mockReturnValueOnce( '/' )
    } );

    const meta = renderer( LayoutAuthenticatedNavigation, {
        directives: {
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
