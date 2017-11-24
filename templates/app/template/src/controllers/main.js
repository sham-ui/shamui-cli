import { DI } from 'sham-ui';
import App from '../widgets/App.sht';

export default function() {
    const app = new App( 'body', 'app', {

    } );

    DI.bind( 'widgets:app', app );
};
