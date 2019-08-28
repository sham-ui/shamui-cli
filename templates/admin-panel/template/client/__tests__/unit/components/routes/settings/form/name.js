import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormName  from '../../../../../../src/components/routes/settings/form/name.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesSettingsFormName, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
