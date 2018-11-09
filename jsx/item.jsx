import React from 'react'

export default class Item extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            data: {}
        };

        if ('data' in props) {
            this.state.data = props.data;
        }

        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeDescription = this.onChangeDescription.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }

    onChangeName(event) {
        var value = event.target.value
        this.setState((state) => {
            state.data.Name = value
            return state
        })
    }

    onChangeDescription(event) {
        var value = event.target.value
        this.setState((state) => {
            state.data.Description = value
            return state
        })
    }

    onSubmit(event) {
        event.preventDefault();
        alert("submitted");
        console.log(this.state.data);

        fetch("/shop/keeper/badges/" + this.state.data.ID, {
            method: this.state.data.ID === "new" ? "POST" : "PUT",
            body: JSON.stringify(this.state.data)
        }).then((response) => response.json()).then((data) => {
            this.setState((state) => {
                state.data = data
                return state
            })
        })
    }

    render() {
        return (
            <div className="item">
                <form onSubmit={this.onSubmit}>
                    <input placeholder="Badge Name"
                        value={this.state.data.Name}
                        onChange={this.onChangeName} />
                    <textarea placeholder="Description"
                        value={this.state.data.Description}
                        onChange={this.onChangeDescription}>
                    </textarea>
                </form>
            </div>
        );
    }
}
