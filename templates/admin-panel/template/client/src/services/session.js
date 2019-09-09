import { DI, inject } from 'sham-ui';

export default class Session {
    @inject router;
    @inject store;
    @inject( 'session:storage' ) data;

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
                this._goToLoginPage();
                return false;
            }
        );
    }

    logout() {
        this.store.logout().then(
            () => {
                this.resetSessionValidation();
                this._goToLoginPage();
            }
        );
    }

    _isRestrictedPage() {
        const page = this.router.lastRouteResolved();
        return undefined !== page && ![
            'signup',
            'login'
        ].includes( page.name );
    }

    _goToLoginPage() {
        this.data.isAuthenticated = false;
        this.data.email = '';
        this.data.name = '';
        if ( this._isRestrictedPage() ) {
            requestAnimationFrame(
                () => this.router.navigate(
                    this.router.generate( 'login', {} )
                )
            );
        }
    }
}
