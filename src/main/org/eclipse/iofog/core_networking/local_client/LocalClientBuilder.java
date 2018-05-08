package org.eclipse.iofog.core_networking.local_client;

import io.netty.channel.Channel;
import org.eclipse.iofog.core_networking.local_client.private_client.PrivateLocalClient;
import org.eclipse.iofog.core_networking.local_client.public_client.PublicLocalClient;
import org.eclipse.iofog.core_networking.main.CoreNetworking;

/**
 * builder class to build private/public local clients
 * <p>
 * Created by saeid on 4/13/16.
 */
public class LocalClientBuilder {
    public static LocalClient build(Channel comSatChannel) {
        if (CoreNetworking.config.getMode().equals("public"))
            return new PublicLocalClient(comSatChannel);
        else
            return new PrivateLocalClient(comSatChannel);
    }
}
