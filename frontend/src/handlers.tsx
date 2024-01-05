export const handleDelete = (url: string, type: string, id: number) => {
    //console.log(String.prototype.concat(url, '/delete/', type, '/', id.toString())
    if (window.confirm("Are you sure you want to delete this?")) {
        fetch(
            String.prototype.concat(url, type, '/', id.toString()), { 
            method: 'DELETE',
            headers: { 
                "Content-Type": "application/json",
                "Accept": "application/json", },
            credentials: 'include',
        }).then((resp) => {
            if (resp.ok) {
                window.location.href = "/account"  
            } else if (resp.status === 401) {
                console.log ("Log in this account to perform this action")
                window.location.href = "/login"
            } else {
                console.log ("error performing action")
                window.location.reload()
            }
        })
    } else {
        console.log("Request cancelled")
    }
}

type patchRequest = {
    input1: string | FormDataEntryValue | null,
    input2: string | FormDataEntryValue | null,
}

export const handlePatch = (url: string, type: string, id: number, data: patchRequest) => {
    if (window.confirm("Are you sure you want to edit this?")) {
        fetch(
            String.prototype.concat(url, '/', type, '/', id.toString()), { 
            method: 'PATCH',
            headers: { 
                "Content-Type": "application/json",
                "Accept": "application/json", },
            credentials: 'include',
            body: JSON.stringify(data)
        }).then((resp) => {
            if (resp.ok) {
                window.location.href = "/account"  
            } else if (resp.status === 401) {
                throw Error("Log in to this account to perform this action")
                //window.location.href = "/login"
            } else {
                throw Error("error performing action")
                //window.location.reload()
            }
        })
    } else {
        console.log("Request cancelled")
    }
}
