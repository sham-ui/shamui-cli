import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import RoutesSettingsPage  from '../../../../../src/components/routes/settings/page.sfc';
import renderer from 'sham-ui-test-helpers';

afterEach( () => {
    DI.resolve( 'session:storage' ).reset();
} );

it( 'renders correctly', () => {
    DI.bind( 'title', {
        change() {}
    } );
    const storage = DI.resolve( 'session:storage' );
    storage.name = 'Test member';
    storage.email = 'test@test.com';
    storage.sessionValidated = true;

    const meta = renderer( RoutesSettingsPage, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
