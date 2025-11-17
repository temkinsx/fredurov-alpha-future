import { Link } from "react-router-dom";
const NotFoundPage = () => {
    return (
        <div>
            <h2>404 - Page Not Found</h2>
            <p>The page you are looking for does not exist.</p>
            <Link to="/">
                <button>Go back home</button>
            </Link>
        </div>
    );
}

export default NotFoundPage;