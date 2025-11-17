import { useState } from "react";

type Callback = () => Promise<void>;

export default function useFetch(callback: Callback) {
    const [isPending, setIsPending] = useState(false);
    const [error, setError] = useState<string>("");

    const fetching = async () => {
        try {
            setIsPending(true);
            await callback();
        } catch (e: any) {
            setError(e.message);
        } finally {
            setIsPending(false);
        }
    }

    return { fetching, isPending, error };
}