import { useState } from 'preact/hooks';
import './HomeTopbar.css'

async function register() {
    const res = await fetch(`/api/register`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: "admin",
            password: "password",
        }),
    });

    const data = await res.json();
    console.log(data)
}

function UserDropDown() {
    const [open, setOpen] = useState(false) 
    register()


    return (
        <div class="user">
            <span onClick={() => setOpen((prev) => !prev)}>User</span>
            {open && (
                <div class="dropdown">
                    <input placeholder="Username" type="text" />
                    <input placeholder="Password" type="password" />
                    <input value="Login" type="submit" />
                  
                </div>
            )}   
        </div> 
    );
}

export default function HomeTopbar() {
    return(
        <nav class="top-navbar">
            <span>A top bar</span>
            <input type="search" placeholder="Search" />
            <UserDropDown />
        </nav>
    );
}