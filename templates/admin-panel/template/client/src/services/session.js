import { DI } from 'sham-ui';
import { inject } from 'sham-ui-macro/babel.macro';

export default class Session {
    @inject router;
    @inject store;
    @inject( 'session:storage' ) data;
    @inject( 'app:storage' ) app;

    constructor() {
        DI.bind( 'session', this );
    }

    login( email, password ) {
        return this.store.login( { email, password } ).then( ( { Email, Name } ) => {
            this.data.email = Email;
            this.data.name = Name;
        } ).then(
            ::this.resetSessionValidation
        ).then(
            () => this.validateSessionPromise
        );
    }

    get validateSessionPromise() {
        if ( undefined === this._validateSessionPromise ) {
            this._validateSessionPromise = this.validateSession();
        }
        return this._validateSessionPromise;
    }

    resetSessionValidation() {
        this._validateSessionPromise = undefined;
    }

    validateSession() {
        this.data.sessionValidated = false;
        return this.store.validSession().then(
            ( { Email, Name } ) => {
                this.data.sessionValidated = true;
                this.data.isAuthenticated = true;
                this.data.email = Email;
                this.data.name = Name;

                // Manual run sync for guaranteed update Layout
                // component before promise resolved
                this.data.sync();
                return true;
            },
            () => {
                this.data.sessionValidated = true;
                this.data.isAuthenticated = false;
                this.data.email = '';
                this.data.name = '';
                this.data.sync();
                return false;
            }
        );
    }

    logout() {
        this.store.logout().then(
            () => {

                // Reset router
                this.app.routerResolved = false;
                this.app.sync();

                // Reset cached session
                this.resetSessionValidation();

                // Go to login page
                this.router.navigate(
                    this.router.generate( 'login' )
                );
            }
        );
    }
}
