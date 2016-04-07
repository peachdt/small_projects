<?php
define('DS', DIRECTORY_SEPARATOR);

function _getExtractSchemaStatement($sqlFileName, $db)
{
    $dumpSchema = 'mysqldump' . ' ';
    $dumpSchema .= '--no-data' . ' ';
    $dumpSchema .= '-u ' . $db['user'] . ' ';
    $dumpSchema .= '-p' . $db['pass'] . ' ';
    $dumpSchema .= '-h ' . $db['host'] . ' ';
    $dumpSchema .= $db['name'] .' > ' . $sqlFileName;

    return $dumpSchema;
}

function _getExtractDataStatement($sqlFileName, $db)
{
    $tables = array(
        'adminnotification_inbox',
        'aw_core_logger',
        'dataflow_batch_export',
        'dataflow_batch_import',
        'log_customer',
        'log_quote',
        'log_summary',
        'log_summary_type',
        'log_url',
        'log_url_info',
        'log_visitor',
        'log_visitor_info',
        'log_visitor_online',
        'index_event',
        'report_event',
        'report_viewed_product_index',
        'report_compared_product_index',
        'catalog_compare_item',
        'catalogindex_aggregation',
        'catalogindex_aggregation_tag',
        'catalogindex_aggregation_to_tag'
    );

    $ignoreTables = ' ';
    foreach($tables as $table) {
        $ignoreTables .= '--ignore-table=' . $db['name'] . '.' . $db['pref'] . $table . ' ';
    }

    $dumpData = 'mysqldump' . ' ';
    $dumpData .= $ignoreTables;
    $dumpData .=  '-u ' . $db['user'] . ' ';
    $dumpData .= '-p' . $db['pass'] . ' ';
    $dumpData .= '-h ' . $db['host'] . ' ';
    $dumpData .= $db['name'] .' >> ' . $sqlFileName;

    return $dumpData;
}

function export_tiny()
{
    $configPath = '.' . DS . 'app' . DS . 'etc' . DS . 'local.xml';
    $xml = simplexml_load_file($configPath, NULL, LIBXML_NOCDATA);

    $db['host'] = $xml->global->resources->default_setup->connection->host;
    $db['name'] = $xml->global->resources->default_setup->connection->dbname;
    $db['user'] = $xml->global->resources->default_setup->connection->username;
    $db['pass'] = $xml->global->resources->default_setup->connection->password;
    $db['pref'] = $xml->global->resources->db->table_prefix;

    echo $db['host'] . PHP_EOL;
    echo $db['name'] . PHP_EOL;
    echo $db['user'] . PHP_EOL;
    echo $db['pass'] . PHP_EOL;
    echo $db['pref'] . PHP_EOL;

//    $sqlFileName =  'var' . DS . $db['name'] . '-' . date('j-m-y-h-i-s') . '.sql';
//
//    //Extract the DB schema
//    $dumpSchema = _getExtractSchemaStatement($sqlFileName, $db);
//    exec($dumpSchema);
//
//    //Extract the DB data
//    $dumpData = _getExtractDataStatement($sqlFileName, $db);
//    exec($dumpData);
}

export_tiny();