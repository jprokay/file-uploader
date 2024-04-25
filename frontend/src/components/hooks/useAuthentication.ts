import { client } from "@/lib/contacts-api";
import { useEffect } from "react";

function useAuthentication() {
	useEffect(() => {
		const userId = document.cookie
			.split(";s")
			.find((row) => row.startsWith("userId="))
			?.split("=")[1];
		if (userId === undefined) {
			client.GET("/auth");
		}
	}, []);
}

export default useAuthentication;
