import axios from 'axios';
import { DI, inject } from 'sham-ui';

export default class Store {
    @inject session;

    constructor() {
        DI.bind( 'store', this );
        const baseURL = PRODUCTION ?
            `${document.location.protocol}//${document.location.host}/api/` :
            'http://localhost:3001/api/';
        this.axios = axios.create( {
            baseURL,
            withCredentials: true
        } );
        this.axios.interceptors.request.use(
            ::this._requestInterceptor
        );
        this.axios.interceptors.response.use(
            ( response ) => response,
            ::this._responseFailInterceptor
        );
        this.tokenLoaded = false;
    }

    _requestInterceptor( request ) {
        if ( !this._isCSRFToken( request ) ) {
            request.headers[ 'X-CSRF-Token' ] = this.token;
        }
        if ( this.session.data.isAuthenticated || (
            this._isCSRFToken( request ) ||
            this._isLoginRequest( request ) ||
            this._isSignupRequest( request ) ||
            this._isValidSession( request ) ||
            this._isLogout( request )
        ) ) {
            return request;
        }
    }

    _responseFailInterceptor( error ) {
        const { request, response } = error;
        if ( this.session.data.isAuthenticated && response && 401 === response.status ) {
            this.session.logout();
            return Promise.reject( error );
        }
        if (
            this._isCSRFToken( request ) ||
            this._isValidSession( request ) ||
            this._isLogout( request )
        ) {
            return Promise.reject( error );
        }
        return Promise.reject(
            this.constructor.extractErrors( error )
        );
    }

    static extractErrors( error ) {
        return error && error.response && error.response.data ?
            error.response.data :
            {};
    }

    static extractData( { data } ) {
        return data;
    }

    csrftoken() {
        return this.axios.get( '/csrftoken' ).then(
            response => response.headers[ 'x-csrf-token' ]
        ).then( token => {
            this.token = token;
            this.tokenLoaded = true;
        } );
    }

    get tokenPromise() {
        if ( undefined === this._tokenPromise ) {
            this._tokenPromise = this.csrftoken();
        }
        return this._tokenPromise;
    }

    validSession() {
        return this.tokenPromise.then(
            () => this.axios.get( '/validsession' )
        ).then(
            this.constructor.extractData
        );
    }

    signUp( data ) {
        return this.axios.post( '/members', data );
    }

    login( data ) {
        return this.axios.post( '/login', data ).then( response => {
            this.token = response.headers[ 'x-csrf-token' ];
            this.tokenLoaded = true;
            return this.constructor.extractData( response );
        } );
    }

    logout() {
        return this.axios.post( '/logout' );
    }

    updateMemberName( data ) {
        return this.axios.put( '/members/name', data );
    }

    updateMemberEmail( data ) {
        return this.axios.put( '/members/email', data );
    }

    _isLoginRequest( request ) {
        return '/login' === request.url && 'post' === request.method;
    }

    _isSignupRequest( request ) {
        return '/members' === request.url && 'post' === request.method;
    }

    _isCSRFToken( request ) {
        return '/csrftoken' === request.url && 'get' === request.method;
    }

    _isValidSession( request ) {
        return '/validsession' === request.url && 'get' === request.method;
    }

    _isLogout( request ) {
        return '/logout' === request.url && 'post' === request.method;
    }
}
