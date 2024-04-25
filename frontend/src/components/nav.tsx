"use client"

import { NavigationMenu, NavigationMenuItem } from "@radix-ui/react-navigation-menu";
import { NavigationMenuLink, NavigationMenuList, navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
import Link from "next/link";
import useAuthentication from "./hooks/useAuthentication";

export function Navigation() {
	useAuthentication()

	return (
		<NavigationMenu>
			<NavigationMenuList>
				<NavigationMenuItem>
					<Link href="/" legacyBehavior passHref>
						<NavigationMenuLink className={navigationMenuTriggerStyle()}>
							Directories
						</NavigationMenuLink>
					</Link>
					<Link href="/contacts" legacyBehavior passHref>
						<NavigationMenuLink className={navigationMenuTriggerStyle()}>
							Contacts
						</NavigationMenuLink>
					</Link>

				</NavigationMenuItem>
			</NavigationMenuList>
		</NavigationMenu>

	)
}
