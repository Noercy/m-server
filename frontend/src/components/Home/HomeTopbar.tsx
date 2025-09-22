import { useEffect, useState } from 'preact/hooks';
import './HomeTopbar.css'

async function register(username: string, password: string) {
    const res = await fetch(`/api/register`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
    });

    const data = await res.json();
    console.log(data)
}

async function getCurrentUser(): Promise<string | null> {
    const res = await fetch("/me", {
        method: "GET",
        credentials: "include"
    });

    if (!res.ok) return null;

    const data = await res.json();
    return data.user || null;
}

async function login(username: string, password: string) {
    const res = await fetch(`/api/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json", },
        body: JSON.stringify({ username, password }),
    });

    const data = await res.json();
    console.log(data)

    if (data) {
        return { success: true, username}
    }
    return { success: false }
}

function UserDropDown() {
    const [open, setOpen] = useState(false) 
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [user, setUser] = useState<string | null>(null);
    
    useEffect(() => {
        (async () => {
            const existingUser = await getCurrentUser();
            if (existingUser) {
                setUser(existingUser)
            }
        })();
    }, []);


    const handleLogin = async (e) => {
        e.preventDefault();

        const formData = new FormData(e.currentTarget);

        setUsername((formData.get('username') as string) || "");    
        setPassword((formData.get('password') as string) || "");       

       const res = await login(username, password)
       if (res.success) {
            setUser(res.username ?? "");
       }
       // await register(username, password)
    }

    return (
        <div class="user">
            <span onClick={() => setOpen((prev) => !prev)}>
                {user ? user : "User"}
            </span>
            {open && (
                <div class="dropdown">
                    {user ? (
                    <div>
                        <p>Welcome, {user}!</p>
                        <button onClick={() => setUser(null)}>Logout</button>
                    </div>
                    ):(
                    <form class="loginform" onSubmit={handleLogin}>
                        <input placeholder="Username" type="text" name="username" autoComplete="on"/>
                        <input placeholder="Password" type="password" name="password" />
                        <button type="submit" name="submit" value="login">Login</button>
                    </form>
                    )}
                </div>
            )}   
        </div> 
    );
}

export default function HomeTopbar() {
    return(
        <nav class="top-navbar">
            <span>A top bar</span>
            <input type="search" name="search" placeholder="Search" autoComplete="on"/>
            <UserDropDown />
        </nav>
    );
}