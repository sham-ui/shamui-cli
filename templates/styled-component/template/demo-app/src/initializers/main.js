import App from '../components/App.sht';

export default function() {
    new App( {
        ID: 'app',
        container: document.querySelector( 'body' )
    } );
}
