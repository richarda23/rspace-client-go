Installation and configuration notes for RSpace 1.58:  April 2019
=================================================================

!! Please note that you should first go through the instructions in the Server Setup guides before attempting to install the RSpace application !!
In particular, you must have set up the database and imported the database schema before proceeding.

These notes give installation and configuration information for the latest  
release of RSpace. 

These instructions assume you are using Tomcat8. If  you are still using Tomcat7 then please 
 replace references to  `tomcat8` with `tomcat7`.

For user-oriented documentation, please see:

https://researchspace.helpdocs.io/

Contents of the release bundle
-------------------------------

The unzipped release should contain:

* RSpace installation notes.
* Server pre-requisite installation notes.
* An example server log file illustrating successful server startup logs.
* The RSpace web application .war file
* A configuration file, deployment.properties
* A folder for 3rd party licenses
* A MySQL script to initialise the database. 

What is covered in this document?

* Installing the Application
* Licensing
* Editing properties
* Configuring Tomcat
* Starting the application
* Login 
* Application updates
* Troubleshooting
* Logging

Installing the Application
--------------------------

The application install is a fairly simple deployment for a Tomcat .war file

First we download the package zip from www.researchspace.com using the username and password issued by ResearchSpace.

    cd ~
    wget --user=USERNAME --ask-password http://researchspace.com/electronic-lab-notebook/media/rspace/rspace-%VERSION%.zip
    unzip rspace-%VERSION%.zip (or version you just downloaded)

Copy config into place

    sudo mkdir /etc/rspace
    sudo cp deployment.properties /etc/rspace/
    sudo cp externalLicenses/license.cxl /etc/rspace/

At the moment the RSpace application will only run within the Tomcat webapp folder /ROOT so this should be cleared and the new war file copied into place.

E.g. on Ubuntu

    sudo rm -rf /var/lib/tomcat8/webapps/ROOT
    sudo cp researchspace-<VERSION>-RELEASE.war /var/lib/tomcat8/webapps/ROOT.war

In addition, currently RSpace needs to be able to create files in the Tomcat home folder.
E.g on Ubuntu:

    sudo chown tomcat8:tomcat8 /var/lib/tomcat8

Configuring Tomcat
------------------

Some variables in Tomcat need to be set for the application to work properly; there are several ways of setting these. However in our example, we will assume no other Tomcat applications are running on the server.

In Ubuntu these settings are set in the following file.

    /etc/default/tomcat8

So we edit this file and set the following;

    sudo vim /etc/default/tomcat8

    JAVA_OPTS="-XX:MaxPermSize=256m -Xms512m -Xmx2048m -XX:+CMSClassUnloadingEnabled\
      -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/media/rspace/logs-audit"
    CATALINA_OPTS="-DpropertyFileDir=file:/etc/rspace/ -Dchemaxon.license.url=/etc/rspace/license.cxl\
     -DRS_FILE_BASE=/PATH/TO/FILESTORAGE -Djava.awt.headless=true\
     -Dliquibase.context=run -Dspring.profiles.active=prod -Djmelody.dir=/media/rspace/jmelody"
    
    JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64

(Adjust java_home variable to be that of the installed java jdk)
PLEASE NOTE: "/PATH/TO/FILESTORAGE" should be the server path to your filestore eg. /data/rspace-filestore 

The `-XX:+HeapDumpOnOutOfMemoryError` and `-XX:HeapDumpPath` are optional arguments that  set a  path to a folder that a heap dump can be written to in the event of an OutOfMemory error.
This is very useful for error diagnosis. A suitable folder would be writable by Tomcat, for example the folder holding the error logs.

The `jmelody.dir` holds a path to a folder that is writable by Tomcat and stores records of CPU/memory usage
 etc for trouble-shooting and monitoring of the server. To create this folder:
 
     cd /etc/rspace
     mkdir jmelody
     chmod 755 jmelody
     sudo chown tomcat8:tomcat8 jmelody
     

NOTE: In CentOS these settings are in the following file.

    /etc/tomcat/tomcat.conf 

### Tomcat performance optimization

The default Tomcat configuration is designed for development usage, not production.
There are several changes you can make to configure Tomcat better for production usage. These are not compulsory but may result in better performance.

####AJP

If you are using AJP protocol in Apache, remember to uncomment the ajp/8009 connector in Tomcat `server.xml` file in `/var/lib/tomcat8/conf`

##### JSP compilationn

In file web.xml in $TOMCAT_HOME/conf folder, add the following init parameters to the 'jsp' servlet.
For more information see the 'JSPs' chapter of Tomcat documentation.

    <init-param>
        <param-name>development</param-name>
        <param-value>false</param-value>
    </init-param>
    <init-param>
        <param-name>trimSpaces</param-name>
        <param-value>true</param-value>
    </init-param>
    <init-param>
        <param-name>genStringAsCharArray</param-name>
        <param-value>true</param-value>
    </init-param>

##### Logging 
Secondly, you can disable console logging: open $TOMCAT_HOME/conf/logging.properties and remove the logger 'java.util.logging.ConsoleHandler' from the list of handlers.
E.g.

    .handlers = 1catalina.org.apache.juli.AsyncFileHandler

By default tomcat makes daily logs of accesses and these build up over time to fill up the default disk, see https://tomcat.apache.org/tomcat-7.0-doc/config/valve.html#Access_Logging.  If you want to limit this behaviour, you can set the org.apache.catalina.valves.AccessLogValve to have the property 'rotatable = "false"', and use the logrotate utility to do some rotation and deletion of older logs.

Licensing
--------

RSpace uses an external license server to monitor and restrict the number of active users  on each installation.
In order to use RSpace effectively, the installation server should have external access over the internet to https://legacy.researchspace.com:8093 to ensure the license can be verified. 

Some functionality requires installation of a 3rd party license:

* Chemistry structure search: Chemaxon server license.

The Chemaxon license is included as part of your RSpace package so no further purchase is necessary.
This license permits the users to conduct searches of specific chemical structures within their entries. The license should be put in /etc/rspace,  as specified in the Tomcat command line setting described above e.g. -Dchemaxon.license.url=/etc/rspace/

* Adding an RSpace license: please see section 'Editing Properties' below.

Editing properties
------------------

There is a property file  called `deployment.properties` in the release which contains some properties that MUST be set for your installation. They appear in that file at the top in the Mandatory Properties section.  Set these mandatory properties before going further.

 By default the `deployment.properties` file can be found in the /etc/rspace directory.
 
### Document previews

A separate application is used to generate  document previews of Office/OpenOffice documents. If the application is not installed, then these previews will not be available to users.

To install, download the application from the RSpace download site:

    wget --user=<username> --password=<password> http://www.researchspace.com/electronic-lab-notebook/media/rspace/aspose-app-VERSION.zip

replacing VERSION with the current release version and using your download credentials in place of <username> and <password>
    
Unzip and follow the installation instructions in the file `Usage.md`.

There are several deployment properties relating to document preview generation. These three are **mandatory** if you want document previewing enabled:

* **aspose.license** 	 	Absolute file path to Aspose license E.g. 

    aspose.license=/etc/rspace/aspose/Aspose-Total-Java.lic

* **aspose.app** 	 	Absolute file path to Aspose standalone document converter E.g.

    aspose.license=/etc/rspace/aspose/aspose-app.jar 	
    
* **aspose.logfile**      Absolute path to Aspose document converter's log file. E.g. 
    
    aspose.logfile=/etc/rspace/aspose/logs.txt

#### Optional changes
    
* **aspose.logLevel** The log level (default is INFO) e.g.
    
    aspose.logLevel=WARN
    
* **aspose.jvmArgs**  Optional jvm args to pass to application. E.g.

    aspose.jvmArgs=-Xmx1024m

### MySQL

You can now override the default MySQL connection settings in `deployment.properties`

The properties are as follows:

* **jdbc.url**  	 The database URL, e.g., jdbc:mysql://localhost:3306/rspace
* **jdbc.username** The MySQL username
* **jdbc.password** The MySQL password

### Email

For small trials, you are welcome to use our email server, but if this is not suitable or you wish to use your own email server, then please edit these settings.
You can edit these in `deployment.properties`

* **mail.default.from** Value of 'from' field of email - i.e., where the email appears to come from 	
* **mail.transport.protocol** Mail protocol - default 'smtp'
* **mail.emailAccount** The mail account of the email server,  e.g., support@xxx.com
* **mail.password** Mail account password
* **mail.emailHost** E.g., auth.smtp.1and1.co.uk
* **mail.from** Overrides mail.default.from
* **mail.port** Default 587
* **mail.replyTo** The reply to address

### External application connections

The following optional properties enable RSpace to connect to Mendeley (if Mendeley integration is enabled):

* **mendeley.id** Numeric application id obtained from Mendeley developer site. Default is unset.
* **mendeley.secret** Secret word obtained from Mendeley developer site. Default is unset.

The following optional properties enable RSpace to connect to OneDrive (if OneDrive integration is enabled)::
* **onedrive.client.id**  application id obtained from OneDrive developer site. Default is unset.
* **onedrive.redirect** Callback URL that is configured on OneDrive developer site for this RSpace. Default is unset.

The following optional properties enable RSpace to connect to Enterprise Box API (if this integration is enabled):
* **box.client.id** Client id of Box App registered for given RSpace instance  
* **box.client.secret** Client secret of Box App registered for given RSpace instance

Properties related to ownCloud configuration:
 
If your organisation uses ownCloud, the following  properties are required to enable RSpace to connect to it:
 
* **owncloud.url** The base URL of the ownCloud, e.g. `owncloud.url=https://owncloud-test.researchspace.com`
* **owncloud.server.name=ownCloud** A display label, e.g. `owncloud.server.name=ownCloud` 
* **owncloud.auth.type** Must be 'oauth' e.g `owncloud.auth.type=oauth`
* **owncloud.client.id** The client ID obtained through registering RSpace as an integration
* **owncloud.secret** The  secret obtained through registering RSpace as an integration

The following optional properties enable RSpace to connect to Orcid  API (if this integration is enabled):
* **orcid.client.id** Client id of Orcid App registered for given RSpace instance  
* **orcid.client.secret** Client secret of Orcid App registered for given RSpace instance

The following optional properties enable RSpace to connect to Github API (if this integration is enabled):
* **github.client.id** Client id of Github registered for given RSpace instance  
* **github.secret** Client secret of Github registered for given RSpace instance

The following optional properties enable RSpace to connect to Figshare API (if this integration is enabled):
* **figshare.id** Client id of Figshare App registered for given RSpace instance  
* **figshare.secret** Client secret of Figshare App registered for given RSpace instance

The following optional properties enable RSpace to connect to Slack (if this integration is enabled):
* **slack.client.id** Client id of Slack App registered for given RSpace instance  
* **slack.secret** Client secret of Slack App registered for given RSpace instance
* **slack.verification.token** Verification token.

The following optional property enable  GoogleHangoutsChat to search RSpace (if this integration is enabled):
*  **ghangoutschat.verification.token**

#### Google's reCAPTCHA on 'Sign up' page (since 1.41)

If you're worried about spam accounts being created on your RSpace instance you can enable captcha field on RSpace 'Sign up' page. We are using [Google's reCAPTCHA](https://www.google.com/recaptcha/intro/index.html) technology and you need a Google account to register the  URL of your instance and to obtain API credentials for captcha.

To enable reCAPTCHA on 'Sign up' page:
1. Go to [https://www.google.com/recaptcha/admin](https://www.google.com/recaptcha/admin). Register a new site for your domain (i.e. myresearchspace.com). Note down the values of 'Site key' and 'Secret key'.
2. Add/update following RSpace deployment properties:
  * user.signup.captcha.enabled=true
  * user.signup.captcha.site.key=<your Site key>
  * user.signup.captcha.secret=<your Secret key>
3. Restart RSpace instance, go to 'Sign up' page. The captcha field should appear, and should be required for signup.


### LDAP connections
These optional settings will enable you to import user data from LDAP:

* **ldap.enabled** enabling user creation with help of LDAP, true/false
* **ldap.authentication.enabled** enabling user authentication through LDAP, true/false
* **ldap.url** ldap server URL, only needed if ldap.enabled is true. E.g. ldaps://kudu.rspace.com
* **ldap.baseSuffix** ldap server url, only needed if ldap.enabled is true. E.g. 'dc=test,dc=kudu,dc=axiope,dc=com'
* **ldap.ignorePartialResultException** if set to 'true' suppresses PartialResultException on search queries

* **ldap.bindQuery.dn** user to use for non-anonymous LDAP bind
* **ldap.bindQuery.password** password to use for non-anonymous LDAP bind
* **ldap.anonymousBind** whether anonymous bind should be used

* **ldap.userSearchQuery.uidField** name of LDAP attribute that is supposed to match the searched username ('uid' by default) 
* **ldap.userSearchQuery.dnField** name of LDAP attribute that is holding full DN of the user, to be used during authentication 

### SSO configuration
Set these properties to configure RSpace to run in SSO mode e.g. for Shibboleth integration. There may be further integration work needed with Apache headers/redirects etc. to get this working. 
* **deployment.standalone** true /false. Set to false to enable SSO integration. Default is true
* **deployment.sso.type=SAML** if single sign-on is configured (if deployment.standalone=false) , this property switches authentication filter to SAML. 
* **deployment.sso.logout.url** the URL to redirect to after logout from RSpace. There is no default. If this is not set, users will get an error page after logging out.
* **user.signup.acceptedDomains**  restricts self sign-up for users in SSO environments.
* **picreateGroupOnSignup.enabled** true/false If users can sign up, then this enables them to sign up as a PI
Only users with a username ending with the accepted domain will be allowed to sign up, and other users will be redirected to an information page. There is no default. E.g., @uni.ac.uk

### External file system configuration

RSpace can link to external file systems via Samba or SFTP protocols. To enable this, set these properties:

* **netfilestores.enabled** - true/false Enables UI and functionality to support file systems.
* **netfilestores.auth.pubKey.passphrase** - If authenticating to a file system by public/private keys, this is the the passphrase securing it. Otherwise this property need not be set.
* **netfilestores.extraSystemProps** - takes a comma-separated list defining additional System Properties that should be set before initialising nfs clients
* **netfilestores.export.enabled**. Activates the option to enable filesystem-linked files to be included in the export
. If set to `true`, an checkbox will appear in the export dialog to optionally include linked files.  Default is `false`. 

### Application behaviour
These optional settings configure some global behaviour of the RSpace application.

* **files.maxUploadSize** The maximum individual file size upload, in bytes. Default is 10Mb
* **max.tiff.conversionSize**, The max file size in bytes at which conversion of TIFF images to .png working images will be attempted. The default is 8192000 (8Mb)
* **ui.bannerImage.path** A URI to a png, gif or jpg that can be used to replace the RSpace logo. E.g. file:/etc/rspace/mylogo.png. Default not set
* **ui.bannerImage.url**. A URL to navigate to after clicking on banner image, when user is logged in. Default is the Workspace, making the banner image act as a "Home" button.
* **ui.bannerImage.loggedOutUrl**. A URL to navigate to after clicking on banner image, when user is not logged in. Default is https://www.researchspace.com/
* **archive.folder.storagetime** Time, in hours, that an exported archive will be available for download before it is considered for physical deletion from the server. Default is 24.
* **archive.minSpaceRequiredToStartMB** Minimum available disk space on temp folder partition required for RSpace to start/continue archive export process. In megabytes, default is 1000 (1 GB). 
* **archive.maxExpandedSizeMB** Maximum size of constructed archive export. In megabytes, default is 10000 (10 GB).
* **importArchiveFromServer.enabled**. true/false.  If true, enhances import of RSpace XML archives by being able to specify a file path on the server. This enables archives stored on the server to be re-imported without tedious download and re-upload through web interface. Default is false.
* **sysadmin.delete.user** true/false If true, Sysadmin user can permanently and irreversibly delete user and all their work from the database. Useful for removing wrongly created user accounts. *Use with caution*. Default is false
* **rs.indexOnstartup**  true/false. Default is true. Controls whether text-data is re-indexed or not at application start-up. Setting to false will result in faster startup times. Default is true
* **sysadmin.limitedIpAddresses.enabled** true/false. If true, will only allow sysadmin login from a whitelist of IP addresses. Default is false.
* **pdf.defaultPageSize** Default PDF export page size . Valid values are UNKNOWN,A4,LETTER. Default is A4.
* **profile.email.editable** true/false. If true, user can edit their email address in their profile page. If false, they cannot. Default is true.
* **profile.firstlastname.editable** true/false. If true, user can edit their display name in their profile page. If false, they cannot. Default is true.
* **pdffont.dir** A directory that contains UTF-8 compliant fonts for PDF export; see 1.53 release notes for more details.
* **example.import.files**.

Custom content can now be provided when a user first logs into RSpace. An RSpace XML zip file can be put on the server,
 in a location specified by this property. The imported content will be put in each user's 'Examples' folder.
e.g.

    example.import.files=file:/etc/rspace/ExampleImport-RSPAC-1789.zip

You can provide multiple files using a comma-separated list. Remember to use an absolute path, and to prefix the path with 'file:'
* **user.signup.signupCode** If your RSpace is on public internet, and you want to enable users to sign themselves up, but don't want other people or bots to signup, you can require your users to enter a code when they sign up. The code is compared case-sensitively with the value of this property. 
This mechanism is activated if code is non-null. E.g. `user.signup.signupCode=mydomain29876`

You could also consider a Captcha or similar mechanism. 




### Other settings
Depending on your installation, you might be advised by ResearchSpace to set the following properties:

* **velocity.ext.dir** the path to a folder where customized Velocity templates are kept. This is only needed if you have customer-specific, bespoke email messages
* **rs.postbatchsignup.emailtemplate** The name of a customised Velocity template to be used for sending to users following batch upload. The template file should be in the folder indicated by **velocity.ext.dir**
* **liquibase.context** If your installation needs to have  specific data pre-loaded into the database, you may need to set this property.
* **licenseserver.poll.cron** Crontab value defining how often to poll license server for updates. Default is every 30 minutes.

### Using an external file store

RSpace can now use a 3rd party file storage service as its backing file store. Initially we are supporting Egnyte file store, which requires the following properties to be set:

 **rs.filestore=EGNYTE** Tells RSpace to expect to use Egnyte as file store
 **rs.ext.filestore.baseURL** The URL of your Egnyte instance, e.g. `rs.ext.filestore.baseURL=https://apprspace.qa-egnyte.com`
 **rs.ext.filestore.root** The top-level folder under which RSpace will save files, e.g. `rs.ext.filestore.root=/Shared/RSpaceFileStore`
 **egnyte.internal.app.client.id** The client ID obtained after registering RSpace as an integration on Egnyte
 
This feature is currently in beta, please ask if you are interested in using this feature.

### Cache configuration
This is an advanced topic and on initial set up will not require configuration. Regular monitoring
 can detect if and when the caches become full.
These  caches can become full when using the default settings. You can see the cache usage in the 
System->Monitoring tab. If either of these caches is at, or near, 100% capacity it is a good idea to resize.
Try doubling the cache size until the cache is no longer full. Resetting these properties requires a restart to take effect. Population of caches depends on usage and may take some time (hours or days) before the  caches are being used fully.
The values are the number of items to store in the cache.
 
* **cache.com.researchspace.model.FileProperty** (default 2500)
* **cache.com.researchspace.model.ImageBlob** (default 1000)
* **cache.com.researchspace.model.User.roles** (default 1000)
* **cache.com.researchspace.model.field.FieldForm** (default 1000)
* **cache.com.researchspace.model.record.RSForm**(default 1000)

### Monitoring

RSpace contains an embedded monitoring capability to monitor CPU usage, memory, server load etc at `/monitoring` . The login credentials are 
   `Username: rs_monitoring, Password: <>`.
   
### Configuring updates

In top level folder:

    git clone https://github.com/ResearchSpace-ELN/rspace-update.git
   
 to install the update script, see the readme file for details
 
### Configuring backups

 If we are responsible for backups:
 
* Create a folder `scripts` in top level folder (usually /home/builder/)
* `cd scripts` and clone Backup code:  `git clone https://github.com/ResearchSpace-ELN/rspace-update.git`
* `cd rspace-update` and follow the instructions in `readme.md`  to set up a backup procedure to AWS S3

If we are *not* responsible for backups this section can be ignored.

Starting the application
------------------------
Starting the application is simple but the command will vary slightly from OS to OS

In Ubuntu we use

    sudo /etc/init.d/tomcat8 start

In Centos we use

    sudo /etc/init.d/tomcat start

However in version 7 onwards we should use the new systemctl commands

    sudo systemctl start tomcat.service
    
Testing the installation
------------------------

Login
--------

There is a hard-coded admin account set up for you, the first time RSpace is run:

Username	Password  	Role

sysadmin1	sysWisc23!	system

Depending upon your configuration, new users can create a new account, in a default user role.
The emails and the name associated with this account can be set in the Profile page
for these accounts (in MyRspace tab)

You *must* change this password immediately for servers that are open to the world (the password is in documentation, and widely known).

System users can view usage activity of the application and 
batch-upload new users via the 'System' tab.

### Browser requirements:

A modern browser that supports HTML5 is essential; these are minimum values :

Chrome v 23.0.1271.95 or later

Firefox v56 or later

Safari version 10.0 or later

Internet Explorer 10 or later

Application updates
------------------

Updating the application is a simple case of redeploying tomcat with the new supplied .war file.

Here are the steps you should follow for updating.

Before starting, please review the update changelogs since you last released updated or installed 
 the application - there may be new configuration settings to apply.

 1. Shutdown Tomcat.  `./bin/shutdown.sh`, or `sudo /etc/init.d/tomcat6 stop`
 2. Backup your database:
    `mysqldump -uecatdev -pecatpwd --database ecat5_demo > RS_oldbackup.sql`
 3. Make a copy of old web application
    `cp -r ./webapps/ROOT   /backup/location/for/webapp`
  and remove old web application from webapps:
    `rm -rf webapps/ROOT*`
 4. Copy researchspace-VERSION.war to webapps
    `cp researchspace-VERSION.war ./webapps/ROOT.war`
 5. Restart Tomcat
    `./bin/startup.sh`, or `sudo /etc/init.d/tomcat8 start` or `sudo service tomcat8 start`.
 6. Check log files and your URL for activity:
    `tail -f ./logs/catalina.out`
If installation is successful, you should see a line like:

    INFO \[-pool-2-thread-1\] 21 Apr 2014 12:34:13,044 - SanityChecker.onAppStartup(52) | Sanity check run: ALL_OK=true

Troubleshooting
---------------

Should you have any trouble installing the application you can contact support by email at  support@researchspace.com

If you have followed the steps details above but are still receiving errors in Tomcat during startup please send the last stack trace to the support email address listed above.

E.g.

    tail -500 $CATALINA_HOME/logs/catalina.out > logs.txt

and attach 'logs.txt' to your email.

### Cleaning the search index

Occasionally, after restarting, there may be a lock on the search index file, causing an error like this:

ERROR | Exception occurred org.apache.lucene.store.LockObtainFailedException: Lock obtain timed out: SimpleFSLock@FTsearchIndices/com.researchspace.model.record.Folder/lucene-

In this case:

* shutdown RSpace
* remove the folder FTSearchIndices via `rm -rf FTSearchIndices`
* restart RSpace

Logging
--------

### Overview 
There are 6 logfiles generated by  RSpace:

* error.log - general logging and errors, that used to go to catalina.out
* SecurityEvents.txt - login/authentication/authorisation events
* RSLogs.txt - audit trail logs
* httpRequests.log - a complete record of URLs and the user/timestamp
* emailErrors.txt - dedicated log for failed emails sent by RSpace.
* SlowRequests.txt - dedicated log for requests that take > 5s to complete on the server.

We log all access to services using log4j. All logs (prior to 1.30) used to be logged to catalina.out as well. By default, these 4 log files are written to Tomcat's home folder.

### Setting the folder location of the log files
If you wish, you can set a deployment property, `logging.dir`, to set a folder where these  log files will be located. E.g. 

    logging.dir=rspacelogs 

will create a subfolder in Tomcat folder called 'rspacelogs' and all logfiles will be written there.
Even better, you can externalise the log folder to be wherever you like - just as long as its writable by Tomcat. E.g.

    logging.dir=/absolute/path/to/logfolder

#### Migrating logs from existing log location
If you are updating an existing RSpace installation, and want to set `logging.dir`, there are two things you must do.
Firstly, move all your audit trail logs to the new folder location (this is so that audit trail search will still work).

E.g. `cd` to Tomcat home folder - there will be some files with names starting with `RSLogs.txt`

    mv RSLogs*  /absolute/path/to/logfolder/ 
    
will transfer the files to the new location ( this *must* be the folder  you set as value of `logging.dir`.

Other logs can be moved across as well, or archived, or deleted, as you see fit. Only RSLog files are consumed by RSpace itself.

Secondly, update the value of the property `sysadmin.errorfile.path` to point to the new error log location. E.g. change 

    sysadmin.errorfile.path=/path/to/tomcat/logs/catalina.out

to

    sysadmin.errorfile.path=/absolute/path/to/logfolder/error.log

### Setting the log level
The current log level is WARN; this generates  logging on application performance. You can, if you wish, set this to INFO (more verbose) or ERROR (less verbose) in 
`$Tomcathome/WEB-INF/classes/log4j.xml`, changing the text in 

    <logger name="com.axiope">
        <level value="WARN"/>
    </logger>
    
 to 'INFO' or 'ERROR' as appropriate.
 
The high level log location is stored in files 'RSLogs.txt'. All activity is logged in log files of maximum 10Mb; by default a maximum of 100000 files will be kept.
 
These logs are used as the raw data for RSpace's audit trail functionality; if they are moved or
deleted then the audit trail search may not return accurate results.
 
End of file