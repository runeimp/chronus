Chronus v0.1.0
==============

Go library and command line tool to help with time parsing, formatting, and accuracy


Features
--------

* [ ] Convert ISO 8160 datetime to
	* [x] UNIX Epoch timestamp
	* [x] SQL datetime string (ISO 8160 based)
	* [x] SQL datetime with timezone string (ISO 8160 based)
	* [x] RFC 3339 DateTime
	* [ ] Format (MM/DD/YYYY, YYYY-MM-DD HH:mm:ss TZ, etc.)
* [ ] Convert format (MM/DD/YYYY, YYYY-MM-DD HH:mm:ss TZ, etc.) date to
	* [ ] UNIX Epoch timestamp
	* [ ] ISO 8160 datetime string
	* [ ] RFC 3339 DateTime
	* [ ] SQL datetime string (ISO 8160 based)
	* [ ] SQL datetime with timezone string (ISO 8160 based)
* [ ] Convert UNIX Epoch timestamp to
	* [ ] ISO 8160 datetime string
	* [x] RFC 3339 DateTime
	* [ ] SQL datetime string (ISO 8160 based)
	* [ ] SQL datetime with timezone string (ISO 8160 based)
	* [ ] Format (MM/DD/YYYY, YYYY-MM-DD HH:mm:ss TZ, etc.)
* [ ] List file times
	* [ ] Access
	* [ ] Modification
	* [ ] Creation
	* [ ] etc.
* [ ] Option to read date from file modification time
* [ ] Parse time in the correct location
	* Automatically based on system settings?
	* When a country code is provided by the `CHRONUS_COUNTRY_CODE` environment variable or `-country-code` command line option
		* [ ] EU
		* [ ] UK
		* [x] US
		* [ ] etc.
* [ ] ____


Example Usage
-------------

### Parsing a Date Time With Timezone but Without Country Code Info

```bash
$ chronus "2021-03-08 16:06:34 MST"
               UNIX Timestamp: 1615219594
             Python Timestamp: 1615219594.000000
                 SQL DateTime: 2021-03-08 16:06:34
   UK Slash Date (DD/MM/YYYY): 08/03/2021
   US Slash Date (MM/DD/YYYY): 03/08/2021
            RFC 3339 DateTime: 2021-03-08T16:06:34Z
                     ISO 8601: 2021-03-08T16:06:34
           ISO 8601 Alternate: 2021-03-08 16:06:34 Z (MST)
ISO 8601 FileSafe 1 w/Seconds: 2021-03-08_160634Z
ISO 8601 FileSafe 2 w/Seconds: 2021-03-08_160634_MST

```

Note the default UTC timezone being referenced (the _Z_ instead of -0700 for MST) even though the numbers are correct. This is a Go issue.


### Parsing a Date Time With Timezone and Country Code Info

```bash
$ chronus -country-code US "2021-03-08 16:06:34 MST"
               UNIX Timestamp: 1615244794
             Python Timestamp: 1615244794.000000
                 SQL DateTime: 2021-03-08 16:06:34
   UK Slash Date (DD/MM/YYYY): 08/03/2021
   US Slash Date (MM/DD/YYYY): 03/08/2021
            RFC 3339 DateTime: 2021-03-08T16:06:34-07:00
                     ISO 8601: 2021-03-08T16:06:34
           ISO 8601 Alternate: 2021-03-08 16:06:34 -0700 (MST)
ISO 8601 FileSafe 1 w/Seconds: 2021-03-08_160634-0700
ISO 8601 FileSafe 2 w/Seconds: 2021-03-08_160634_MST

```

This does everything correctly because we supplied the right country code reference which is necessary as many countries often use the same abbreviated versions for timezones, i.e.; EST is use locally in most countries that have an Eastern Standard Time. Such as Australian Eastern Standard Time (AEST) is typically referred to as EST in Australia. So most timezone abbreviations can not be trusted.

Note that `CHRONUS_COUNTRY_CODE="US" chronus "2021-03-08 16:06:34 MST"` would have the same effect as the example above.

At this time "US" and "USA" are the only country codes supported. As there is interest, from myself or others, more country code support will be added. To have Chronus always default to US set the environment variable `CHRONUS_COUNTRY_CODE=US` or `CHRONUS_COUNTRY_CODE=USA`. Others countries will also be referenced via their standard Alpha-2 and Alpha-3 codes when their support is added.


Articles & Reference
--------------------

* [3.3. Date and Time Specification - RFC 5322][]
	* Obsoletes [3.3. Date and Time Specification - RFC 2822][]
		* Obsoletes [5. DATE AND TIME SPECIFICATION - RFC 822][]
* [List of TZ Database Time Zones - Wikipedia][]
* [Sources for time zone and daylight saving time data - IANA][]





[3.3. Date and Time Specification - RFC 5322]: https://datatracker.ietf.org/doc/html/rfc5322#page-14
[3.3. Date and Time Specification - RFC 2822]: https://datatracker.ietf.org/doc/html/rfc2822#page-14
[5. DATE AND TIME SPECIFICATION - RFC 822]: https://datatracker.ietf.org/doc/html/rfc822#section-5
[List of TZ Database Time Zones - Wikipedia]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
[Sources for time zone and daylight saving time data - IANA]: https://data.iana.org/time-zones/tz-link.html

