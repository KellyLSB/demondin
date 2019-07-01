export default class StateErrors {
	constructor(props) {
		super(props);
		
		this.onError = this.onError.bind(this);
		this.hasError = this.hasError.bind(this);
		this.getError = this.getError.bind(this);
	}

	onError(name, msg = null) {
		this.setState((state) => {
			if (msg === false) {
				delete state.errors[name];
			} else {
				state.errors[name] = msg;
			}
			
			return state;
		} );
	}
	
	hasError(name = null) {
		if (name === null) {
			return Object.keys(this.state.errors).length > 0;
		}
		
		return this.state.errors.hasOwnProperty(name);
	}
	
	getErrors() {
		return Object.values(this.state.errors);
	}
}
