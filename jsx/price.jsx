import React from 'react'

export default class Price extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            data: 'data' in props ? props.data : {}
        };
    }

    render() {
        return (
            <form>
                <input placeholder="$0.00"
                    onChange={this.onChangePrice}
                    value={this.state.new.price} />
                <input placeholder="DatePicker"
                    onChange={this.onChangeBeforeDate}
                    value={this.state.data.before_date} />
                <input placeholder="DatePicker"
                    onChange={this.onChangeAfterDate}
                    value={this.state.data.after_date} />
            </form>
        )
    }
}
