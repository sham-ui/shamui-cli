import * as directives from 'sham-ui-directives';
import RoutesSettingsParagraph  from '../../../../../src/components/routes/settings/paragraph.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesSettingsParagraph, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'default onUpdate options', () => {
    const meta = renderer( RoutesSettingsParagraph, {
        directives,
        form: compile`<button data-test-dummy-button :onclick=\{{onUpdateSuccess}}>Click me!</button>`
    } );
    meta.component.container.querySelector( '.icon-pencil' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
    meta.component.container.querySelector( '[data-test-dummy-button]' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
} );
