import { DI } from 'sham-ui';
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

it( 'display errors', async() => {
    expect.assertions( 2 );

    const updateMock = jest.fn();
    DI.bind( 'store', {
        updateMemberName: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormName, {
        directives: {
            ...directives
        }
    } );

    const formData = {
        newName: 'Johny Smithy'
    };
    const { component } = meta;
    component.container.querySelector( '[name="name"]' ).value = formData.newName;
    component.container.querySelector( '[type="submit"]' ).click();
    component.container.querySelector( '[data-test-modal] [data-test-ok-button]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
